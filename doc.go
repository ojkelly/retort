/*
Retort is a reactive terminal user interface framework for golang.

Inspired by React, the API is somewhat similar, but due to langauge differences they are
not the same thing.

Components

An app built with retort is composed of components.

Hooks

Retort uses hooks to provide the functionality to make your Components interactive and responsive
to user input and data.

There are a few built in hooks, which can also be used to create custom hooks.

UseState: use this to keep track of, and change state that is local to a component

UseEffect: use this to do something (like setState) in a goroutine, for example fetch data

UseScreen: use this to access the screen object directly. You probably wont need this, but it's there
if you do for example if you want to create a new ScreenElement.

UseQuit: use this to exit the application.

Why

As stated by the inspiration for this package "Declarative views make your code more predictable and
easier to debug.Spew". The original author (Owen Kelly) has years of experience building complex websites
with React, and wanted a similar reactive/declarative tool for terminal user interfaces in golang.

The biggest reason though, is state management.

When you build an interactive user interface, the biggest challenge is always state management. The
model that a reactive framework like retort allows, is one of the simplest ways to solve the
state management problem. Much moreso than an imperitive user interface library.

About the Name

retort: to answer (an argument) by a counter argument

Terminals usually have arguments. Don't think about it too much.


Examples

Below are some simple examples of how to use retort

*/
package retort // import "retort.dev"
