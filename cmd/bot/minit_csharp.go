// +build csharp

package main

// Initialize any modules that are in their own packages and require the csharp tag

import (
	_ "github.com/pajbot/pajbot2/pkg/modules/message_height_limit"
)
