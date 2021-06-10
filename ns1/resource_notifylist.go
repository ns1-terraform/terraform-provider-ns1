package ns1

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/monitor"
)

func notifyListResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notifications": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"config": {
							Type:     schema.TypeMap,
							Required: true,
						},
					},
				},
			},
		},
		Create: NotifyListCreate,
		Read:   NotifyListRead,
		Update: NotifyListUpdate,
		Delete: NotifyListDelete,
	}
}

func notifyListToResourceData(d *schema.ResourceData, nl *monitor.NotifyList) error {
	d.SetId(nl.ID)
	d.Set("name", nl.Name)

	if len(nl.Notifications) > 0 {
		notifications := make([]map[string]interface{}, len(nl.Notifications))
		for i, n := range nl.Notifications {
			ni := make(map[string]interface{})
			ni["type"] = n.Type
			if n.Config != nil {
				ni["config"] = n.Config
			}
			notifications[i] = ni
		}
		d.Set("notifications", notifications)
	}
	return nil
}

func resourceDataToNotifyList(nl *monitor.NotifyList, d *schema.ResourceData) error {
	nl.ID = d.Id()

	if rawNotifications := d.Get("notifications").([]interface{}); len(rawNotifications) > 0 {
		ns := make([]*monitor.Notification, len(rawNotifications))
		for i, notificationRaw := range rawNotifications {
			ni := notificationRaw.(map[string]interface{})
			config := ni["config"].(map[string]interface{})

			switch ni["type"].(string) {
			case "user":
				ns[i] = monitor.NewUserNotification(config["user"].(string))
			case "email":
				ns[i] = monitor.NewEmailNotification(config["email"].(string))
			case "datafeed":
				ns[i] = monitor.NewFeedNotification(config["sourceid"].(string))
			case "webhook":
				ns[i] = monitor.NewWebNotification(config["url"].(string))
			case "pagerduty":
				ns[i] = monitor.NewPagerDutyNotification(config["service_key"].(string))
			case "hipchat":
				ns[i] = monitor.NewHipChatNotification(config["token"].(string), config["room"].(string))
			case "slack":
				ns[i] = monitor.NewSlackNotification(config["url"].(string), config["username"].(string), config["channel"].(string))
			default:
				return fmt.Errorf("%s is not a valid notifier type", ni["type"])
			}
		}
		nl.Notifications = ns
	}
	return nil
}

// NotifyListCreate creates an ns1 notifylist
func NotifyListCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	nl := monitor.NewNotifyList(d.Get("name").(string))

	if err := resourceDataToNotifyList(nl, d); err != nil {
		return err
	}

	if resp, err := client.Notifications.Create(nl); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return notifyListToResourceData(d, nl)
}

// NotifyListRead fetches info for the given notifylist from ns1
func NotifyListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	nl, resp, err := client.Notifications.Get(d.Id())
	if err != nil {
		if err == ns1.ErrListMissing {
			log.Printf("[DEBUG] NS1 notify list (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return ConvertToNs1Error(resp, err)
	}

	return notifyListToResourceData(d, nl)
}

// NotifyListDelete deletes the given notifylist from ns1
func NotifyListDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	resp, err := client.Notifications.Delete(d.Id())
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

// NotifyListUpdate updates the notifylist with given parameters
func NotifyListUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	nl := monitor.NewNotifyList(d.Get("name").(string))

	if err := resourceDataToNotifyList(nl, d); err != nil {
		return err
	}

	if resp, err := client.Notifications.Update(nl); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return notifyListToResourceData(d, nl)
}
