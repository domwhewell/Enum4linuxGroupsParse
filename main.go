package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var recurse bool
var file_name string
var username string
var group_name string

var help = fmt.Sprint("This script will recursively lookup groups or users in an active directory environment from the enum4linux groups output.\n- For example 'user1' is a member of \x1b[31m'IT Support'\x1b[0m which is itself a member of \x1b[31m'Domain Admins'\x1b[0m making 'user1' effectively a domain admin.\nThis script will colour code the output so \x1b[31m'Groups'\x1b[0m are always highlighted as \x1b[31mred\x1b[0m to make them easier to make out.\n")

func groups(groupsRaw, username, groupName string, recurse bool) {
	fmt.Println(help)

	// Obtain a list of group names
	group_lines := regexp.MustCompile(`\x1b?(?:\^\[)?(?:\[[0-9;]*m)?group:\[(.+?)\] rid:\[.+?\]`).FindAllStringSubmatch(groupsRaw, -1)
	groups := []string{}
	for _, match := range group_lines {
		groups = append(groups, match[1])
	}

	if groupName != "" {
		// If the domain is included in the group name, remove it.
		if strings.Contains(groupName, "\\") {
			groupName = strings.Split(groupName, "\\")[1]
		}

		groupMembers(groupName, groups, groupsRaw)
	} else if username != "" {
		// If the domain is included in the username, remove it.
		if strings.Contains(username, "\\") {
			username = strings.Split(username, "\\")[1]
		}

		membersOfGroup(username, groupsRaw, groups)
	} else {
		fmt.Println("\x1b[31m[!]\x1b[0m Error you did not specify a group_name or username")
	}
}

func groupMembers(groupName string, groups []string, groupsRaw string) {
	fmt.Println("\x1b[32m[!]\x1b[0m Members of group: \x1b[31m" + groupName + "\x1b[0m")
	var name string
	var group_names []string
	members_of_group := regexp.MustCompile(`\x1b?(?:\^\[)?(?:\[[0-9;]*m)?Group: \x1b?(?:\^\[)?(?:\[[0-9;]*m)?\'?`+groupName+`\'? \(RID\: \d+\) has member\: (.+)`).FindAllStringSubmatch(groupsRaw, -1)
	for _, member := range members_of_group {
		if strings.Contains(member[1], "\\") {
			name = strings.Split(member[1], "\\")[1]
		} else {
			name = member[1]
		}
		if contains(groups, name) {
			if !contains(group_names, name) {
				group_names = append(group_names, name)
			}
		} else {
			fmt.Println(name)
		}
		if recurse {
			for _, group := range group_names {
				fmt.Println("\x1b[32m[!]\x1b[0m\x1b[31m " + group + "\x1b[0m is a member of \x1b[31m" + group_name + "\x1b[0m")
				groupMembers(group, groups, groupsRaw)
			}
		}
	}
}

func membersOfGroup(username, groupsRaw string, groups []string) {
	groupsUserIsAMemberOf := regexp.MustCompile(`(?m)\x1b?(?:\^\[)?(?:\[[0-9;]*m)?Group\:? \x1b?(?:\^\[)?(?:\[[0-9;]*m)?\'?(.+?)\'? \(RID\: \d+\) has member\: .+\\`+username+`$`).FindAllStringSubmatch(groupsRaw, -1)
	for _, member_group := range groupsUserIsAMemberOf {
		// Check if this variable is a group name or a username so we can colour code them differently
		var member_of string
		if contains(groups, member_group[1]) {
			member_of = fmt.Sprint("\x1b[31m" + member_group[1] + "\x1b[0m")
		} else {
			member_of = member_group[1]
		}

		// Check if this variable is a group name or a username so we can colour code them differently
		var member string
		if contains(groups, username) {
			member = fmt.Sprint("\x1b[31m" + username + "\x1b[0m")
		} else {
			member = username
		}
		fmt.Println(member + " is a member of " + member_of)
		if recurse {
			membersOfGroup(member_group[1], groupsRaw, groups)
		}
	}
}

func contains(x []string, s string) bool {
	for _, v := range x {
		if v == s {
			return true
		}
	}
	return false
}

func main() {
	// grab some arguments from the command line
	flag.BoolVar(&recurse, "recurse", true, "Disable recursive search")
	flag.StringVar(&file_name, "file", "", "Set the `enum4linux -G` file to read from")
	flag.StringVar(&username, "username", "", "Set the username to search for")
	flag.StringVar(&group_name, "group_name", "", "Set the groupname to search for")
	flag.Parse()

	if file_name != "" {
		data, err := ioutil.ReadFile(file_name)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Convert the buffer to a string.
		groups_raw := string(data)

		groups(groups_raw, username, group_name, recurse)
	} else {
		fmt.Println(help)
		flag.Usage()
		fmt.Println()
		fmt.Println("\x1b[31m[!]\x1b[0m Error you must specify the 'enum4linux -G' output file to load.")
	}
}
