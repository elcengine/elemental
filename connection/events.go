package e_connection

import elemental "github.com/elcengine/elemental/core"

// Add a listener to a connection
//
// Deprecated: Use elemental.OnConnectionEvent instead.
var On = elemental.OnConnectionEvent

// Remove a listener from a connection
//
// Deprecated: Use elemental.RemoveConnectionEvent instead.
var Off = elemental.RemoveConnectionEvent
