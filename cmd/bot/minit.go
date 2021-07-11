package main

// Initialize any modules that are in their own packages

import (
	_ "github.com/pajbot/pajbot2/pkg/modules/commands" // xd

	_ "github.com/pajbot/pajbot2/pkg/modules/tusecommands" // xd

	_ "github.com/pajbot/pajbot2/pkg/modules/punisher" // xd

	_ "github.com/pajbot/pajbot2/pkg/modules/system"

	_ "github.com/pajbot/pajbot2/pkg/modules/giveaway"

	_ "github.com/pajbot/pajbot2/pkg/modules/link_filter"

	_ "github.com/pajbot/pajbot2/pkg/modules/twitter"
)
