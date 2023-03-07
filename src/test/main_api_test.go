package main

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"user-check/model"
	"user-check/utils"
	"user-check/utils/logger"
	"os"
	"strconv"
	"testing"
)

var (
	ctx context.Context
	// api url
	apiUrl string
	//endpoint
	resourceUserCheck string
	resourceUserCount string
	resourceRealUser  string
)

func init() {
	ctx = context.Background()
	logger.Init(ctx, true)

	// initialize test endpoint env var
	apiUrl = utils.EnvOrDefault("TEST_ENDPOINT", "http://localhost:8080/")
	// api endpoints
	resourceUserCount = "/api/v1/usercount"
	resourceUserCheck = "/api/v1/usercheck"
	// this is real user
	resourceRealUser = "martih"

}

// TestMain start testing framework
func TestMain(m *testing.M) {
	{
		log.Println("Starting group-license tests")
		// prepare the tests
		exitVal := m.Run()
		// clean up code here
		os.Exit(exitVal)
	}
}

// TestUserExists test suite to check user functionality
func TestUserExists(t *testing.T) {
	Convey(`Feature: test user check functionality endpoint`, t, func() {
		// no user
		ValidateUserCheckEndpointNoUser(t)
		// random invalid user
		ValidateUserCheckEndpointInvalidUser(t)
		// real user
		ValidateUserCheckEndpointRealUser(t)
	})
}

// TestUserCount test suite to count user
func TestUserCount(t *testing.T) {
	Convey(`Feature: test user count functionality endpoint`, t, func() {
		UserCount(t)
	})
}

// ValidateUserCheckEndpointNoUser call user check without isid param
func ValidateUserCheckEndpointNoUser(t *testing.T) {
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resourceUserCheck
	urlStr := u.String()

	client, _ := funcHttpClientGet(urlStr)

	resp, _ := client.Get(urlStr)
	//data, _ := ioutil.ReadAll(resp.Body)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	//log.Printf("returned:%s", data)
	//log.Printf("status code: %v", resp.StatusCode)

	// check response code 404
	Convey("Check response code", func() {
		Convey("Check response code", func() {
			So(resp.StatusCode, CheckResponse, 404)
		})
	})

}

// ValidateUserCheckEndpointInvalidUser check if user exists in ldap sec group using invalid user
func ValidateUserCheckEndpointInvalidUser(t *testing.T) {
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resourceUserCheck
	urlStr := u.String() + "/" + RandomChar(10)

	client, _ := funcHttpClientGet(urlStr)

	resp, _ := client.Get(urlStr)
	data, _ := ioutil.ReadAll(resp.Body)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	svc1 := model.JSONSuccessResult{}
	jsonErr := json.Unmarshal(data, &svc1)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	Convey("Check user doesnt exists and status code", func() {
		Convey("Check response code", func() {
			So(resp.StatusCode, CheckResponse, 200)
		})

		Convey("Check is user is not in ldap", func() {
			So(fmt.Sprintf("%v", svc1.Data), ShouldEqual, "false")
		})
	})

}

// ValidateUserCheckEndpointRealUser check if user exists in ldap sec group using real user
func ValidateUserCheckEndpointRealUser(t *testing.T) {
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resourceUserCheck
	urlStr := u.String() + "/" + resourceRealUser

	client, _ := funcHttpClientGet(urlStr)

	resp, _ := client.Get(urlStr)
	data, _ := ioutil.ReadAll(resp.Body)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	svc1 := model.JSONSuccessResult{}
	jsonErr := json.Unmarshal(data, &svc1)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//log.Printf("response returned from api:%s", data)
	//log.Printf("status code: %v", resp.StatusCode)
	//log.Printf("data returned: %v", svc1.Data)

	Convey("Check user exists and status code", func() {
		Convey("Check response code for 200", func() {
			// check response code 200
			So(resp.StatusCode, CheckResponse, 200)
		})
		Convey("Check if user is really in ldap", func() {
			So(fmt.Sprintf("%v", svc1.Data), ShouldEqual, "true")
		})
	})

}

// UserCountEndpoint check if count users endpoint returns integer value
func UserCount(t *testing.T) {
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resourceUserCount
	urlStr := u.String()

	client, _ := funcHttpClientGet(urlStr)

	resp, _ := client.Get(urlStr)
	data, _ := ioutil.ReadAll(resp.Body)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	svc1 := model.JSONSuccessResult{}
	jsonErr := json.Unmarshal(data, &svc1)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//log.Printf("returned:%s", data)
	//log.Printf("status code: %v", resp.StatusCode)

	Convey("Check user count and status code", func() {
		Convey("Check response code for 200", func() {
			So(resp.StatusCode, CheckResponse, 200)
		})
		Convey("Check if returned count value is really integer", func() {
			mynumber, err := strconv.Atoi(fmt.Sprintf("%v", svc1.Data))
			if err != nil {
				log.Fatal("something really really bad happened converting str to int, really bad :(")
			}
			// just check type integer, nothing fancy
			So(mynumber, ShouldHaveSameTypeAs, 0)
		})
	})
}

// CheckResponse check response value
func CheckResponse(actual interface{}, expected ...interface{}) string {
	if actual == expected[0] {
		return ""
	}
	return "Response code returned by test is different than what you expected"
}
