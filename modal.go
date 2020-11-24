// Modal example - How to respond to a slash command with an interactive modal and parse the response
// The flow of this example:
// 1. User trigers your app with a slash command (e.g. /modaltest) that will send a request to http://URL/slash and respond with a request to open a modal
// 2. User fills out fields first and last name in modal and hits submit
// 3. This will send a request to http://URL/modal and send a greeting message to the user

// Note: Within your slack app you will need to enable and provide a URL for "Interactivity & Shortcuts" and "Slash Commands"
// Note: Be sure to update YOUR_SIGNING_SECRET_HERE and YOUR_TOKEN_HERE
// You can use ngrok to test this example: https://api.slack.com/tutorials/tunneling-with-ngrok
// Helpful slack documentation to learn more: https://api.slack.com/interactivity/handling

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/slack-go/slack"
)

func generateModalRequest(channelId string) slack.ModalViewRequest {
	// Create a ModalViewRequest with a header and two inputs
	titleText := slack.NewTextBlockObject("plain_text", "Squawk", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Cancel", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "Squawk your way to easier deployments", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	memberOptions := createOptionBlockObjects([]string{"all", "staging", "hk", "production"}, false)
	targetServersText := slack.NewTextBlockObject(slack.PlainTextType, "Target Server List", false, false)
	targetServersOption := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, nil, "targetServers", memberOptions...)
	targetServersBlock := slack.NewInputBlock("targetServers", targetServersText, targetServersOption)

	memberOptions = createOptionBlockObjects([]string{"all", "proxy.conf", "enterprise.conf", "bidder.conf", "partners.conf", "provider.conf", "vitals.conf", "traffic.conf", "region.yml", "proto.json", "supply.conf"}, false)
	ixConfsText := slack.NewTextBlockObject(slack.PlainTextType, "Ix Confs List", false, false)
	ixConfsOption := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, nil, "ixConfs", memberOptions...)
	ixConfsBlock := slack.NewInputBlock("ixConfs", ixConfsText, ixConfsOption)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			targetServersBlock,
			ixConfsBlock,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

// This was taken from the slash example
// https://github.com/slack-go/slack/blob/master/examples/slash/slash.go
func verifySigningSecret(r *http.Request) error {
	signingSecret := signingsecret
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("Received body: %s\n", body)

	// Need to use r.Body again when unmarshalling SlashCommand and InteractionCallback
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	verifier.Write(body)
	if err = verifier.Ensure(); err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func handleSlash(w http.ResponseWriter, r *http.Request) {

	err := verifySigningSecret(r)
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	api := slack.New(slacktoken)

	if s.UserName == "joshua.gross" || s.UserName == "adam.wong" || s.UserName == "joshua.lo" || s.UserName == "michael.tardibuono" || s.UserName == "ryan.pietrow" {
		modalRequest := generateModalRequest(s.ChannelID)
		_, err = api.OpenView(s.TriggerID, modalRequest)
		if err != nil {
			fmt.Printf("Error opening view: %s", err)
		}
	} else {
		msg := fmt.Sprintf("User Not authorized!!!")

		_, _, err = api.PostMessage(s.ChannelID,
			slack.MsgOptionText(msg, false),
			slack.MsgOptionAttachments())

	}

	/*
	   // Not leaving options in, only onoe responsee
	   	switch s.Text {
	   	case "humboldttest":
	   		api := slack.New(slacktoken)
	   		modalRequest := generateModalRequest()
	   		_, err = api.OpenView(s.TriggerID, modalRequest)
	   		if err != nil {
	   			fmt.Printf("Error opening view: %s", err)
	   		}
	   	default:
	   		w.WriteHeader(http.StatusInternalServerError)
	   		return
	   	}
	*/
}

func handleModal(w http.ResponseWriter, r *http.Request) (string, string) {

	err := verifySigningSecret(r)
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return "", ""
	}

	var i slack.InteractionCallback
	err = json.Unmarshal([]byte(r.FormValue("payload")), &i)
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return "", ""
	}

	// Note there might be a better way to get this info, but I figured this structure out from looking at the json response
	targetServers := i.View.State.Values["targetServers"]["targetServers"].SelectedOption.Value
	ixConfs := i.View.State.Values["ixConfs"]["ixConfs"].SelectedOption.Value

	return targetServers, ixConfs
	//msg := fmt.Sprintf("Hello %s %s, nice to meet you!", firstName, lastName)
	/*
		api := slack.New(slacktoken)
		_, _, err = api.PostMessage(i.User.ID,
			slack.MsgOptionText(msg, false),
			slack.MsgOptionAttachments())
		if err != nil {
			fmt.Printf(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	*/
}

func createOptionBlockObjects(options []string, users bool) []*slack.OptionBlockObject {
	optionBlockObjects := make([]*slack.OptionBlockObject, 0, len(options))
	var text string
	for _, o := range options {
		if users {
			text = fmt.Sprintf("<@%s>", o)
		} else {
			text = o
		}
		optionText := slack.NewTextBlockObject(slack.PlainTextType, text, false, false)
		optionBlockObjects = append(optionBlockObjects, slack.NewOptionBlockObject(o, optionText, nil))
	}
	return optionBlockObjects
}
