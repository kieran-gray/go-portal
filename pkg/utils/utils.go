package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"slices"
	"sort"
	"strings"
	"time"

	dt "github.com/kieran-gray/go-portal/pkg/types"
)

func ToMap(data interface{}) map[string]reflect.Value {
	// Assuming keys and values are same length
	m := make(map[string]reflect.Value)
	keys := reflect.TypeOf(data)
	values := reflect.ValueOf(data)
	for i := 0; i < keys.NumField(); i++ {
		value := values.Field(i)
		m[keys.Field(i).Name] = value
	}
	return m
}

func SortedByPriority(environments []dt.Environment) []dt.Environment {
	sort.Slice(environments, func(i, j int) bool {
		return environments[i].Priority < environments[j].Priority
	})
	return environments
}

func HasLogs(service dt.ServiceDetails) bool {
	for _, environment := range service.Environments {
		if environment.LogsUrl != "" {
			return true
		}
	}
	return false
}

func GetHighestPriorityUrl(service dt.ServiceDetails) string {
	return SortedByPriority(service.Environments)[0].Url
}

func GenerateServiceId(service dt.Service) string {
	return strings.Replace(service.Metadata.Name, " ", "", -1)
}

func GetFormattedTimeSince(dateString string) string {
	plural := func(diff int64, str string) string {
		if diff > 1 {
			return fmt.Sprintf("%d %ss Ago", int(diff), str)
		}
		return fmt.Sprintf("%d %s Ago", int(diff), str)
	}
	formatTime := func(timeCreated time.Time) string {
		periods := map[string]int64{
			"year":   365 * 24 * 60 * 60,
			"month":  30 * 24 * 60 * 60,
			"week":   7 * 24 * 60 * 60,
			"day":    24 * 60 * 60,
			"hour":   60 * 60,
			"minute": 60,
		}
		var diff = time.Now().Unix() - timeCreated.Unix()
		switch {
		case diff > periods["year"]:
			return "Over a Year Ago"
		case diff > periods["month"]:
			return plural(diff/periods["month"], " Month")
		case diff > periods["week"]:
			return plural(diff/periods["week"], " Week")
		case diff > periods["day"]:
			return plural(diff/periods["day"], " Day")
		case diff > periods["hour"]:
			return plural(diff/periods["hour"], " Hour")
		case diff > periods["minute"]:
			return plural(diff/periods["minute"], " Minute")
		default:
			return "Just now"
		}
	}
	timestamp, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		log.Print(err.Error())
		return ""
	}
	return formatTime(timestamp)
}

func GetFavourites(r *http.Request) []string {
	cookie, err := r.Cookie("favourites")
	if err != nil {
		return []string{}
	}
	return strings.Split(cookie.Value, " ")
}

func GetDisplayServices(services []dt.Service, favourites []string) dt.DisplayServices {
	displayServices := dt.DisplayServices{Favourites: []dt.Service{}, Services: []dt.Service{}}
	for _, service := range services {
		if slices.Contains(favourites, GenerateServiceId(service)) {
			displayServices.Favourites = append(displayServices.Favourites, service)
		} else {
			displayServices.Services = append(displayServices.Services, service)
		}
	}
	return displayServices
}

func ParseServicesFromRequest(r *http.Request) (dt.ServicesFile, error) {
	var services dt.ServicesFile
	err := json.NewDecoder(r.Body).Decode(&services)
	if err != nil {
		return services, err
	}
	return services, nil
}

func AddService(services []dt.Service) []dt.Service {
	service := dt.Service{
		Metadata: dt.ServiceMetadata{
			Name:    fmt.Sprintf("New Service %d", len(services)),
			Aliases: "",
			DevOnly: false,
		},
		Ui: dt.ServiceDetails{
			Description:   "",
			RepositoryUrl: "",
		},
		Api: dt.ServiceDetails{
			Description:   "",
			RepositoryUrl: "",
		},
	}
	newServices := []dt.Service{service}
	return append(newServices, services...)
}

func AddEnvironment(services []dt.Service, serviceId string, serviceType string) []dt.Service {
	newEnvironment := dt.Environment{
		Name:     "New Environment",
		Priority: 0,
		Url:      "",
		LogsUrl:  "",
	}
	for i := 0; i < len(services); i++ {
		if GenerateServiceId(services[i]) == serviceId {
			if serviceType == "ui" {
				services[i].Ui.Environments = append(services[i].Ui.Environments, newEnvironment)
			} else {
				services[i].Api.Environments = append(services[i].Api.Environments, newEnvironment)
			}
		}
	}
	return services
}

func GetRepoType(repositoryUrl string) string {
	if strings.Contains(repositoryUrl, "gitlab") {
		return "gitlab"
	} else {
		return "github"
	}
}
