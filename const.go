package main

var modal = `{
        "title": {
                "type": "plain_text",
                "text": "Modal Title"
        },
        "submit": {
                "type": "plain_text",
                "text": "Submit"
        },
        "blocks": [
                {
                        "type": "input",
                        "element": {
                                "type": "plain_text_input",
                                "action_id": "sl_input",
                                "placeholder": {
                                        "type": "plain_text",
                                        "text": "Placeholder text for single-line input"
                                }
                        },
                        "label": {
                                "type": "plain_text",
                                "text": "Label"
                        },
                        "hint": {
                                "type": "plain_text",
                                "text": "Hint text"
                        }
                },
                {
                        "type": "input",
                        "element": {
                                "type": "plain_text_input",
                                "action_id": "ml_input",
                                "multiline": true,
                                "placeholder": {
                                        "type": "plain_text",
                                        "text": "Placeholder text for multi-line input"
                                }
                        },
                        "label": {
                                "type": "plain_text",
                                "text": "Label"
                        },
                        "hint": {
                                "type": "plain_text",
                                "text": "Hint text"
                        }
                }
        ],
        "type": "modal"
}`

var serverModal = `{
	"title": {
		"type": "plain_text",
		"text": "My App",
		"emoji": true
	},
	"submit": {
		"type": "plain_text",
		"text": "Submit",
		"emoji": true
	},
	"type": "modal",
	"close": {
		"type": "plain_text",
		"text": "Cancel",
		"emoji": true
	},
	"blocks": [
		{
			"type": "section",
			"block_id": "inventory-block",
			"text": {
				"type": "mrkdwn",
				"text": "Pick servers to push to"
			},
			"accessory": {
				"action_id": "inventory",
				"type": "multi_static_select",
				"placeholder": {
					"type": "plain_text",
					"text": "Select items"
				},
				"options": [
					{
						"text": {
							"type": "plain_text",
							"text": "all"
						},
						"value": "all-server"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "hk"
						},
						"value": "hk-server"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "a321"
						},
						"value": "a321-server"
					}
				]
			}
		},
		{
			"type": "section",
			"block_id": "confs-block",
			"text": {
				"type": "mrkdwn",
				"text": "Pick ix-confs to push"
			},
			"accessory": {
				"action_id": "ix-confs",
				"type": "multi_static_select",
				"placeholder": {
					"type": "plain_text",
					"text": "Select items"
				},
				"options": [
					{
						"text": {
							"type": "plain_text",
							"text": "all"
						},
						"value": "all-conf"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "benchmarking"
						},
						"value": "benchmarking-conf"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "ib"
						},
						"value": "ib-conf"
					}
				]
			}
		}
	]
}`
