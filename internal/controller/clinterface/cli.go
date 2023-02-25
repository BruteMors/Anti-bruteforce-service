package clinterface

import (
	"Anti-bruteforce-service/internal/controller/httpapi/handlers"
	"Anti-bruteforce-service/internal/domain/entity"
	"Anti-bruteforce-service/internal/domain/service"
	"fmt"
	"github.com/c-bata/go-prompt"
	"log"
	"os"
	"strings"
)

var suggestions = []prompt.Suggest{
	{Text: "blacklist add [ip_address] [mask]", Description: "Add ip net to blacklist"},
	{Text: "blacklist remove [ip_address] [mask]", Description: "Remove ip net to blacklist"},
	{Text: "blacklist get", Description: "Get ip list from blacklist"},
	{Text: "whitelist add [ip_address] [mask]", Description: "Add ip net to whitelist"},
	{Text: "whitelist remove [ip_address] [mask]", Description: "Remove ip net to whitelist"},
	{Text: "whitelist get", Description: "Get ip list from whitelist"},
	{Text: "bucket remove [login] [ip_address]", Description: "Remove login and ip address from bucket"},
	{Text: "help", Description: "Display list of commands"},
	{Text: "exit", Description: "Exit anti bruteforce app"},
}

type CommandLineInterface struct {
	serviceAuth *service.Authorization
	serviceWL   *service.WhiteList
	serviceBL   *service.BlackList
}

func New(serviceAuth *service.Authorization, serviceWL *service.WhiteList, serviceBL *service.BlackList) *CommandLineInterface {
	return &CommandLineInterface{serviceAuth: serviceAuth, serviceWL: serviceWL, serviceBL: serviceBL}
}

func (c *CommandLineInterface) Run(ch chan os.Signal) {
	executor := prompt.Executor(func(s string) {
		s = strings.TrimSpace(s)
		setCommand := strings.Split(s, " ")
		switch setCommand[0] {
		case "blacklist":
			switch setCommand[1] {
			case "add":
				if len(setCommand) != 4 {
					break
				}
				c.addIpToBl(entity.IpNetwork{
					Ip:   setCommand[2],
					Mask: setCommand[3],
				})
			case "remove":
				if len(setCommand) != 4 {
					break
				}
				c.removeIpToBl(entity.IpNetwork{
					Ip:   setCommand[2],
					Mask: setCommand[3],
				})
			case "get":
				c.getIpListFromBl()
			default:
				fmt.Println("unknown command")
			}

		case "whitelist":
			switch setCommand[1] {
			case "add":
				if len(setCommand) != 4 {
					break
				}
				c.addIpToWl(entity.IpNetwork{
					Ip:   setCommand[2],
					Mask: setCommand[3],
				})
			case "remove":
				if len(setCommand) != 4 {
					break
				}
				c.removeIpToWl(entity.IpNetwork{
					Ip:   setCommand[2],
					Mask: setCommand[3],
				})
			case "get":
				c.getIpListFromWl()
			default:
				fmt.Println("unknown command")
			}

		case "bucket":
			if len(setCommand) != 4 {
				break
			}
			if setCommand[1] == "reset" {
				c.resetBucket(entity.Request{
					Login:    setCommand[2],
					Password: "",
					Ip:       setCommand[3],
				})
			} else {
				fmt.Println("unknown command")
			}
		case "exit":
			ch <- os.Interrupt
			return
		case "help":
			for _, suggestion := range suggestions {
				fmt.Println("Command:", suggestion.Text, "Description:", suggestion.Description)
			}

		default:
			fmt.Println("unknown command")
		}

	})
	completer := prompt.Completer(func(in prompt.Document) []prompt.Suggest {
		w := in.GetWordBeforeCursor()
		if w == "" {
			return []prompt.Suggest{}
		}
		return prompt.FilterHasPrefix(suggestions, w, true)
	})
	defer func() {
		if a := recover(); a != nil {
			log.Println("Command line interface not available. Please run container with tty mode")
		}
	}()
	prompt.New(executor, completer).Run()
}

func (c *CommandLineInterface) addIpToBl(ipNet entity.IpNetwork) {
	isValidateIP := handlers.ValidateIP(ipNet)
	if !isValidateIP {
		fmt.Println("not valid ip")
		return
	}
	err := c.serviceBL.AddIP(ipNet)
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}
	fmt.Printf("add address: %v to blacklist \n", ipNet)
}

func (c *CommandLineInterface) removeIpToBl(ipNet entity.IpNetwork) {
	isValidateIP := handlers.ValidateIP(ipNet)
	if !isValidateIP {
		fmt.Println("not valid ip")
		return
	}
	err := c.serviceBL.RemoveIP(ipNet)
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}
	fmt.Printf("remove address: %v from blacklist \n", ipNet)
}

func (c *CommandLineInterface) getIpListFromBl() {
	list, err := c.serviceBL.GetIPList()
	if err != nil {
		return
	}
	for _, network := range list {
		fmt.Println(network)
	}
}

func (c *CommandLineInterface) addIpToWl(ipNet entity.IpNetwork) {
	isValidateIP := handlers.ValidateIP(ipNet)
	if !isValidateIP {
		fmt.Println("not valid ip")
		return
	}
	err := c.serviceWL.AddIP(ipNet)
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}
	fmt.Printf("add address: %v to whitelist \n", ipNet)
}

func (c *CommandLineInterface) removeIpToWl(ipNet entity.IpNetwork) {
	isValidateIP := handlers.ValidateIP(ipNet)
	if !isValidateIP {
		fmt.Println("not valid ip")
		return
	}
	err := c.serviceWL.RemoveIP(ipNet)
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}
	fmt.Printf("remove address: %v from whitelist \n", ipNet)
}

func (c *CommandLineInterface) getIpListFromWl() {
	list, err := c.serviceWL.GetIPList()
	if err != nil {
		return
	}
	for _, network := range list {
		fmt.Println(network)
	}
}

func (c *CommandLineInterface) resetBucket(request entity.Request) {
	isValidateReq := handlers.ValidateRequest(request)
	if !isValidateReq {
		fmt.Println("not valid ip")
		return
	}
	isResetIp := c.serviceAuth.ResetIpBucket(request.Ip)
	if !isResetIp {
		fmt.Printf("ip address: %v not find\n", request.Ip)
	} else {
		fmt.Printf("ip address: %v has been reseted\n", request.Ip)
	}

	isResetLogin := c.serviceAuth.ResetLoginBucket(request.Ip)
	if !isResetLogin {
		fmt.Printf("login: %v not find\n", request.Ip)
	} else {
		fmt.Printf("ip address: %v has been reseted\n", request.Ip)
	}
}
