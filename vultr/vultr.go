package vultr

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"git.beanjs.com/wangjun/beanjs-v2ray/config"
	"github.com/vultr/govultr/v3"
	"golang.org/x/oauth2"
)

var client *govultr.Client

func init() {
	ctx := context.Background()

	auth2 := &oauth2.Config{}
	ts := auth2.TokenSource(ctx, &oauth2.Token{AccessToken: config.Configuration.VultrKey})

	client = govultr.NewClient(oauth2.NewClient(ctx, ts))
}

func instanceRemove() error {
	ctx := context.Background()

	ins, _, _, err := client.Instance.List(ctx, &govultr.ListOptions{
		Label: config.Configuration.VultrLabel,
	})

	if err != nil {
		return err
	}

	if len(ins) == 1 {
		log.Printf("vultr instance remove: %v", ins[0].ID)
		if err := client.Instance.Delete(ctx, ins[0].ID); err != nil {
			return err
		}
	}

	return nil
}

func instanceCreate(fwg string) (*govultr.Instance, error) {
	ctx := context.Background()

	ins, _, err := client.Instance.Create(ctx, &govultr.InstanceCreateReq{
		Label:           config.Configuration.VultrLabel,
		Region:          "ewr",
		Plan:            "vc2-1c-0.5gb",
		OsID:            2076,
		FirewallGroupID: fwg,
	})

	if err != nil {
		return nil, err
	}

	log.Printf("vultr instance create: %v", ins.ID)

	for {
		time.Sleep(time.Duration(config.Configuration.Timeout) * time.Second)

		wIns, _, err := client.Instance.Get(ctx, ins.ID)
		if err != nil {
			return nil, err
		}

		log.Printf("vultr instance status: %v", wIns.Status)
		if wIns.Status == "active" {
			ins.MainIP = wIns.MainIP
			break
		}
	}

	return ins, nil
}

func firewallRemove() error {
	ctx := context.Background()
	desc := fmt.Sprintf("%v-firewall", config.Configuration.VultrLabel)

	fws, _, _, err := client.FirewallGroup.List(ctx, &govultr.ListOptions{})

	if err != nil {
		return err
	}

	if len(fws) > 0 {
		for _, v := range fws {
			if v.Description == desc {
				log.Printf("vultr firewall remove: %v", v.ID)
				if err := client.FirewallGroup.Delete(ctx, v.ID); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func firewallCreate() (*govultr.FirewallGroup, error) {

	ctx := context.Background()

	fwg, _, err := client.FirewallGroup.Create(ctx, &govultr.FirewallGroupReq{
		Description: fmt.Sprintf("%v-firewall", config.Configuration.VultrLabel),
	})

	if err != nil {
		return nil, err
	}

	log.Printf("vultr firewall create: %v", fwg.ID)

	client.FirewallRule.Create(ctx, fwg.ID, &govultr.FirewallRuleReq{
		IPType:     "v4",
		Protocol:   "TCP",
		Subnet:     "0.0.0.0",
		SubnetSize: 0,
		Port:       strconv.Itoa(config.Configuration.VultrPort),
		Notes:      "v2ray",
	})

	client.FirewallRule.Create(ctx, fwg.ID, &govultr.FirewallRuleReq{
		IPType:     "v4",
		Protocol:   "TCP",
		Subnet:     "0.0.0.0",
		SubnetSize: 0,
		Port:       "22",
		Notes:      "ssh",
	})

	return fwg, nil
}

func ShowInfo() error {
	ctx := context.Background()

	act, _, err := client.Account.Get(ctx)
	if err != nil {
		return err
	}

	log.Printf("vultr account name: %v", act.Name)
	log.Printf("vultr account email: %v", act.Email)
	log.Printf("vultr account credit: %v", -(act.PendingCharges + act.Balance))
	return nil
}

func Instance(fwg string) (*govultr.Instance, error) {
	if err := instanceRemove(); err != nil {
		return nil, err
	}

	return instanceCreate(fwg)
}

func Firewall() (*govultr.FirewallGroup, error) {
	if err := firewallRemove(); err != nil {
		return nil, err
	}

	return firewallCreate()
}
