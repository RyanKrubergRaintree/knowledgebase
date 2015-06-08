package cmds

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/raintreeinc/knowledgebase/kb"
)

func init() {
	Register(Command{
		Name: "setup-groups",
		Desc: "Setup groups and communities based on configuration file.",
		Run:  SetupGroups,
	})
}

func SetupGroups(DB kb.Database, fs *flag.FlagSet, args []string) {
	conf := fs.String("conf", "", "configuration file to be loaded")
	fs.Parse(args)

	if *conf == "" {
		fmt.Println("Configuration file should be supplied")
		fs.Usage()
		return
	}

	data, err := ioutil.ReadFile(*conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	var definition struct {
		Groups    []kb.Group
		Community []struct {
			GroupID  kb.Slug
			MemberID kb.Slug
			Rights   kb.Rights
		}
	}

	if err := json.Unmarshal(data, &definition); err != nil {
		fmt.Println(err)
	}

	context := DB.Context("admin")
	for _, group := range definition.Groups {
		err := context.Groups().Create(group)
		if err == nil {
			fmt.Printf("%15v: added %v\n", group.ID, group.Name)
		} else {
			fmt.Printf("%15v: failed to add %v: %v\n", group.ID, group.Name, err)
		}
	}

	for _, community := range definition.Community {
		err := context.Access().CommunityAdd(community.GroupID, community.MemberID, community.Rights)
		if err == nil {
			fmt.Printf("%15v: added community %v as %v\n", community.GroupID, community.MemberID, community.Rights)
		} else {
			fmt.Printf("%15v: failed adding %v as %v: %v\n", community.GroupID, community.MemberID, community.Rights, err)
		}
	}
}
