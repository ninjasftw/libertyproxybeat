package beater

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	//"regexp"
	"strconv"
	"strings"
	"time"
	"crypto/tls"
	"crypto/x509"
    "encoding/json"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

const (
	managerJmxproxy = "/IBMJMXConnectorREST/mbeans/"
)

func (bt *Libertyproxybeat) GetJMX(u url.URL) error {
	for i := 0; i < len(bt.Beans); i++ {

		//for each bean create a new map to store all the attributes
        
		logp.Debug("GetJMX", "name :", bt.Beans[i].Name, "size of attr:", len(bt.Beans[i].Attributes))
		var attslice = make([]string, len(bt.Beans[i].Attributes))
		

		for j := 0; j < len(bt.Beans[i].Attributes); j++ {
			//store names into a slice
			logp.Debug("GetJMX","attr name:", bt.Beans[i].Attributes[j].Name)
			attslice[j]=bt.Beans[i].Attributes[j].Name

		}

		err := bt.GetJMXObject(u, bt.Beans[i].Name, attslice, bt.CAFile)
		if err != nil {
			logp.Err("Error requesting JMX: %v", err)
		}

	}
	return nil
}

func (bt *Libertyproxybeat) GetJMXObject(u url.URL, name string, attributes []string, CAFile []uint8) error {

	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool()}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
    var ParsedUrl *url.URL

    if len(CAFile) > 0 {
		ok := tlsConfig.RootCAs.AppendCertsFromPEM(CAFile)
		if !ok {
		    logp.Err("Unable to load CA file")
			panic("Couldn't load PEM data")
		}
    }

	//client := &http.Client{}
	client := &http.Client{Transport: transport}

	ParsedUrl, err := url.Parse(u.String())
    if err != nil {
		logp.Err("Unable to parse URL String")
		panic(err)
    }

    ParsedUrl.Path += managerJmxproxy + url.QueryEscape(name) + "/attributes"
    parameters := url.Values{}

	for _, attr := range attributes {
		parameters.Add("attribute", attr)
	}


	ParsedUrl.RawQuery = parameters.Encode()
	logp.Debug(selector, "Requesting JMX: %s", ParsedUrl.String())  

	req, err := http.NewRequest("GET", ParsedUrl.String(), nil)

	if bt.auth {
		req.SetBasicAuth(bt.username, bt.password)
	}
	res, err := client.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("HTTP %s", res.Status)
	}

	scanner := bufio.NewScanner(res.Body)
	scanner.Scan()

	attrmap, err := GetJMXValue(scanner.Text())
	if err != nil {
		return err
	}

	event := common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "lib_jmx",
		"fields": bt.beatConfig.Libertyproxybeat.Fields,
		"bean": common.MapStr{
			"name":      name,
			"attributes":     attrmap,
			"hostname":  u.Host,
		},
	}
	
	bt.events.PublishEvent(event)
	//logp.Info("Event: %+v", event)

	return nil
}

func GetJMXValue(responseBody string) (map[string]string, error) {

    responseMap := make(map[string]string)

	if strings.HasPrefix(responseBody, "Error") {
		return responseMap, errors.New(responseBody)
	}
	logp.Debug("GetJMXValue", "Response Body: %s", responseBody)  

    var dat []map[string]interface{}


	if err := json.Unmarshal([]byte(responseBody), &dat); err != nil {
        panic(err)
    }

	var beanitem map[string]interface{}

	//TODO: only handles a single bean currently
    for key := range dat {
		beanitem = dat[key]["value"].(map[string]interface{})
        logp.Debug("GetJMXValue", "record endpoint: %s", beanitem["value"])


		// the Liberty JMX api returns a java type float which needs to converted before being returnd
		if beanitem["type"] == "java.lang.Double" {
		    logp.Debug("GetJMXValue", "Double type detected, modifying to float64")
		    modstr := strings.Replace(beanitem["value"].(string), "E", "e+", 1)
		    floatcvt, err := strconv.ParseFloat(modstr, 32)
		    if err == nil {
				responseMap[dat[key]["name"].(string)] = strconv.FormatFloat(floatcvt, 'f', 0, 64)
		    }
		 } else  {
			responseMap[dat[key]["name"].(string)] = beanitem["value"].(string)
         } 

    }

     

	//return beanitem["value"].(string), nil
	return responseMap, nil
}
