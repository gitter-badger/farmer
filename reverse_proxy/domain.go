package reverse_proxy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"errors"
	"github.com/farmer-project/farmer/db"
	"github.com/farmer-project/farmer/farmer"
)

func AddDomain(box *farmer.Box, url string, port string) error {
	if err := checkDomain(url); err != nil {
		return err
	}

	if err := box.AddDomain(url, port); err != nil {
		return err
	}

	if err := ConfigureDomains(box); err != nil {
		return err
	}

	return db.DB.Save(&box).Error
}

func DeleteDomains(box *farmer.Box) error {
	for _, domain := range box.Domains {
		if err := DeleteDomain(box, domain.Url); err != nil {
			return err
		}
	}

	return nil
}

func DeleteDomain(box *farmer.Box, url string) error {
	domain := &farmer.Domain{
		BoxId: box.ID,
		Url:   url,
	}

	if err := db.DB.Find(domain).Error; err != nil {
		return err
	}

	if err := box.DeleteDomain(url); err != nil {
		return err
	}

	os.Remove(configFile(domain.Url, domain.Port))

	return db.DB.Delete(domain).Error
}

func ConfigureDomains(box *farmer.Box) error {
	for _, domain := range box.Domains {
		if err := setReverseProxyConfig(box.IP, domain.Url, domain.Port); err != nil {
			return err
		}
	}

	return nil
}

func checkDomain(url string) error {
	domain := &farmer.Domain{}

	err := db.DB.Where("url = ?", url).Find(domain).Error
	if err != nil {
		// Could not find a domain so domain is available and free to use.
		if err.Error() == "record not found" {
			return nil
		}
		return err
	}

	if domain.Url == url {
		box, err := farmer.FindBoxById(domain.BoxId)
		if err != nil {
			return err
		}

		return errors.New("Domain '" + url + "' is already taken by '" + box.Name + "' box.")
	}

	return nil
}

func setReverseProxyConfig(ip string, url string, port string) error {
	t := template.New("reverse_proxy")
	absPath, err := filepath.Abs("reverse_proxy/nginx.cfg")
	if err != nil {
		return err
	}

	reProxyFile, err := ioutil.ReadFile(absPath)
	if err != nil {
		return err
	}

	t, err = t.Parse(string(reProxyFile))
	if err != nil {
		return err
	}

	type reverseProxy struct {
		Domain      string
		ContainerIP string
		Port        string
	}

	filePath := configFile(url, port)

	os.Remove(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	return t.Execute(file, &reverseProxy{
		Domain:      url,
		ContainerIP: ip,
		Port:        port,
	})
}

func configFile(url string, port string) string {
	reverseProxyConfigFileLocation := os.Getenv("FARMER_REVERSE_PROXY_LOCATION")
	return reverseProxyConfigFileLocation + "/" + url + "." + port + ".conf"
}
