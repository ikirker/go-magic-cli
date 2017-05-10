# go-magic-cli

This is a work-in-progress Go reimplementation of [@slackhq](https://github.com/slackhq/)'s [magic-cli](https://github.com/slackhq/magic-cli).

The idea is to make a collection of related commands in separate scripts look like one command in one script. I originally started reimplementing it because on the HPC clusters we work on we change up environment variables a lot, so especially things that depend on scripting languages can break if something becomes not available or changes version from the one the author expected. Maybe I should have done it in bash instead of Go, though.
