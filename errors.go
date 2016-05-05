/*
 * Copyright 2015 Fabrício Godoy
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package eval

import (
	"fmt"
)

const (
	// Info defines a informational event.
	Info = 0

	// Warn defines a event which indicates wrong parameters.
	Warn = 1

	// Error defines a failed event.
	Error = 2

	// Fatal defines a unrecoverable event.
	Fatal = 3
)

// A NotRunningError represents an error when a Environment function is called
// before it is running.
type NotRunningError string

// Error returns string representation of current instance error.
func (e NotRunningError) Error() string {
	return fmt.Sprintf("The environment '%s' is not running", string(e))
}

// A ErrUser represents an event that needs user attention.
type ErrUser struct {
	Severity int
	Message  string
}

// Error returns string representation of current instance error.
func (e ErrUser) Error() string {
	return e.Message
}
