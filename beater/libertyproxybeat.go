package beater

import (
	"fmt"
	"net/url"
	"time"
	"io/ioutil"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"github.com/ninjasftw/libertyproxybeat/config"
)

const selector = "libertyproxybeat"

type Libertyproxybeat struct {
	beatConfig *config.Config
	done       chan struct{}
	period     time.Duration
	urls       []*url.URL
	auth       bool
	username   string
	password   string
    CAFile     []uint8
	Beans      []Bean
	events     publisher.Client
    fields     map[string]string
}

type Bean struct {
	Name       string
	Attributes []config.Attribute
	Keys       []string
}

// Creates beater
func New() *Libertyproxybeat {
	return &Libertyproxybeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (bt *Libertyproxybeat) Config(b *beat.Beat) error {

	// Load beater beatConfig
	err := cfgfile.Read(&bt.beatConfig, "")
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	return nil
}

func (bt *Libertyproxybeat) Setup(b *beat.Beat) error {

	bt.events = b.Publisher.Connect()

    if len(bt.beatConfig.Libertyproxybeat.Fields) > 0 {
       fmt.Println("YAY")
    }

	// Setting default period if not set
	if bt.beatConfig.Libertyproxybeat.Period == "" {
		bt.beatConfig.Libertyproxybeat.Period = "1s"
	}

	var err error
	bt.period, err = time.ParseDuration(bt.beatConfig.Libertyproxybeat.Period)
	if err != nil {
		return err
	}

	//define default URL if none provided
	var urlConfig []string
	if bt.beatConfig.Libertyproxybeat.URLs != nil {
		urlConfig = bt.beatConfig.Libertyproxybeat.URLs
	} else {
		urlConfig = []string{"http://127.0.0.1:8888"}
	}

	bt.urls = make([]*url.URL, len(urlConfig))
	for i := 0; i < len(urlConfig); i++ {
		u, err := url.Parse(urlConfig[i])
		if err != nil {
			logp.Err("Invalid JMX url: %v", err)
			return err
		}
		bt.urls[i] = u
	}

    if bt.beatConfig.Libertyproxybeat.Ssl.Cafile != "" {
        logp.Info("CAFile IS set.")
		pemdata, err := ioutil.ReadFile(bt.beatConfig.Libertyproxybeat.Ssl.Cafile)
		if err != nil {
			logp.Debug("Config", "Failed to load CA file")
			panic(err)
		}
		bt.CAFile = pemdata
    } else {
        logp.Info("CAFile IS NOT set.")
    }
    

	if bt.beatConfig.Libertyproxybeat.Authentication.Username == "" || bt.beatConfig.Libertyproxybeat.Authentication.Password == "" {
		logp.Err("Username or password IS NOT set.")
		bt.auth = false
	} else {
		bt.username = bt.beatConfig.Libertyproxybeat.Authentication.Username
		bt.password = bt.beatConfig.Libertyproxybeat.Authentication.Password
		bt.auth = true
		logp.Info("Username and password IS set.")
	}

	bt.Beans = make([]Bean, len(bt.beatConfig.Libertyproxybeat.Beans))
	if bt.beatConfig.Libertyproxybeat.Beans == nil {
		logp.Err("No beans are configured set.")
		//TODO: default values (HeapMemory)?
	} else {
		for i := 0; i < len(bt.beatConfig.Libertyproxybeat.Beans); i++ {
			bt.Beans[i].Name = bt.beatConfig.Libertyproxybeat.Beans[i].Name
			bt.Beans[i].Attributes = bt.beatConfig.Libertyproxybeat.Beans[i].Attributes
			bt.Beans[i].Keys = bt.beatConfig.Libertyproxybeat.Beans[i].Keys

			logp.Debug(selector, "Bean name: %s", bt.beatConfig.Libertyproxybeat.Beans[i].Name)
			for j := 0; j < len(bt.beatConfig.Libertyproxybeat.Beans[i].Attributes); j++ {
				logp.Debug(selector, "\tBean attribute: %s", bt.beatConfig.Libertyproxybeat.Beans[i].Attributes[j].Name)
				if len(bt.beatConfig.Libertyproxybeat.Beans[i].Attributes[j].Keys) > 0 {
					for k := 0; k < len(bt.beatConfig.Libertyproxybeat.Beans[i].Attributes[j].Keys); k++ {
						logp.Debug(selector, "\t\tAttribute key: %s", bt.beatConfig.Libertyproxybeat.Beans[i].Attributes[j].Keys[k])
					}
				}
			}
			for k := 0; k < len(bt.beatConfig.Libertyproxybeat.Beans[i].Keys); k++ {
				logp.Debug(selector, "\tBean key: %s", bt.beatConfig.Libertyproxybeat.Beans[i].Keys[k])
			}
		}
	}

	return nil
}

func (bt *Libertyproxybeat) Run(b *beat.Beat) error {
	logp.Info("Libertyproxybeat is running! Hit CTRL-C to stop it.")

	//for each url
	for _, u := range bt.urls {

		go func(u *url.URL) {
			ticker := time.NewTicker(bt.period)
			defer ticker.Stop()

			for {
				select {
				case <-bt.done:
					goto GotoFinish
				case <-ticker.C:
				}

				err := bt.GetJMX(*u)
				if err != nil {
					logp.Err("Error while getttig JMX: %v", err)
				}
			}
		GotoFinish:
		}(u)
	}

	<-bt.done
	return nil
}

func (bt *Libertyproxybeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Libertyproxybeat) Stop() {
	close(bt.done)
}
