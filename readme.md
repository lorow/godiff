# What is this?

Godiff is basically a simple tool for debugging two requests side by side.
It's a very specific tool for a very specific use case - testing requests and comparing the output.

[todo write more about it]

# How to debug?

Since we're using bubbletea as the TUI framework, debugging is a bit awkward. The general idea is to build
the binary with all debug flags enabled, run it and then attach to the running process with a debugger

go build -gcflags="all=-N -l" -o godiff.exe && ./godiff.exe


# architecture 

// here's the new idea:
// inspired by the posting.sh thing
// I wanna get rid of the bottom command bar
// instead it's gonna live in a popup - figure out how to make popup with bubble tea
// this popup will have a search input that's gonna go through all of the commands
// and dispatch the proper one once selected.

// I'd also want to have a jump mechanism where by inputting jump button
// we go into jump state, in that state we can cancel or jump focus to something

// landing page view
// progress on that part:
// - UI mostly done - need to rethink styles
// - I need to handle the filtering command
// - I need to handle shortcut inputs - done
//|--------------------------------------------------------------------------------------------------------------------|
//|  GoDiff - 1.0.0                                                                                                    |
//|                                                                                                                    |
//|   | Search for a project                                                                                       |   |
//|                                                                                                                    |
//|   - Projects - 2 -----------------------------------------------------------------------------------------------   |
//|   |                                  				                                                    	             |   |
//|   |  Project name                                 	                                                    	     |   |
//|   |    Short project description                   	                                                           |   |
//|   |                                  				                                                                   |   |
//|   |  Project name                                 	                                                   	       |   |
//|   |    Short project description                   	                                                   	       |   |
//|   |                                  				                                                                   |   |
//|   --------------------------------------------------------------------------------------------------------------   |
//|                                                                                                                    |
//| up/down select project ^n new project ^o commands enter - load                                                     |
//|--------------------------------------------------------------------------------------------------------------------|

// new project wizard - will be used to only setup a new project
//|--------------------------------------------------------------------------------------------------------------------|
//|  GoDiff - 1.0.0                                                                                                    |
//|                                                                                                                    |
//|   | New Project Title                                                                                          |   |
//|                                                                                                                    |
//|   - Description  -----------------------------------------------------------------------------------------------   |
//|   |                                  				                                                    	             |   |
//|   |                                                                                                            |   |
//|   |                                  				                                                                   |   |
//|   --------------------------------------------------------------------------------------------------------------   |
//|                                                                                         | cancel |  | continue |   |
//| up/down change field esc go back                                                                                   |
//|--------------------------------------------------------------------------------------------------------------------|

// project page view - single editor
//|--------------------------------------------------------------------------------------------------------------------|
//|  GoDiff - 1.0.0                                                                                                    |
//|                                                                                                                    |
//|   |------------------------------------------------------------------------------------------------------------|   |
//|   | GET | http://some-service.dev/                                                                  |          |   |
//|   |------------------------------------------------------------------------------------------------------------|   |
//|                                                                                                                    |
//|   |------------------------------------------------------------------------------------------------------------|   |
//|   |    Short project description                                                                               |   |
//|   |                                                                                                            |   |
//|   |  Project name                                                                                              |   |
//|   |    Short project description                                                                               |   |
//|   |                                                                                                            |   |
//|   --------------------------------------------------------------------------------------------------------------   |
//|                                                                                                                    |
//| <-/up/down/-> change focus ^s save ^o commands ^p jump i edit                                                      |
//|--------------------------------------------------------------------------------------------------------------------|

// project page view - double editor
//|--------------------------------------------------------------------------------------------------------------------|
//|  GoDiff - 1.0.0                                                                                                    |
//|                                                                                                                    |
//|   |-----------------------------------------------------|    |-------------------------------------------------|   |
//|   | GET | http://some-service.dev/                 |  	|    | GET | http://some-service.dev/            |     |   |
//|   |-----------------------------------------------------|    |-------------------------------------------------|   |
//|                                                                                                                    |
//|   |-----------------------------------------------------|    |-------------------------------------------------|   |
//|   |    Short project description                        |    |    Short project description                    |   |
//|   |                                                     |    |                                                 |   |
//|   |  Project name                                       |    |  Project name                                   |   |
//|   |    Short project description                        |    |    Short project description                    |   |
//|   |                                                     |    |                                                 |   |
//|   -------------------------------------------------------    ---------------------------------------------------   |
//|                                                                                                                    |
//| ^c exit ^s save ^o commands ^p jump i edit enter send                                                              |
//|--------------------------------------------------------------------------------------------------------------------|

// command popup
//|--------------------------------------------------------------------------------------------------------------------|
//|                                                                                                                    |
//|                           |------------------------------------------------------------|                           |
//|                           | Search for command                                         |                           |
//|                           |------------------------------------------------------------|                           |
//|                           |  Some command                                              |                           |
//|                           |  Some command explanation                                  |                           |
//|                           |                                                            |                           |
//|                           |  Some command                                              |                           |
//|                           |  Some command explanation                                  |                           |
//|                           |                                                            |                           |
//|                           |  Some command                                              |                           |
//|                           |  Some command explanation                                  |                           |
//|                           |                                                            |                           |
//|                           |------------------------------------------------------------|                           |
//|                                                                                                                    |
//|                                                                                                                    |
//|--------------------------------------------------------------------------------------------------------------------|