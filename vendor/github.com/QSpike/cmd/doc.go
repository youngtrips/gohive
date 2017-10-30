/*
Python's cmd package in golang.

Sample:
    type MyCmd struct {}
    func (this MyCmd) Help() {
        println("available commands:")
        println("list foo")
    }

    func (this MyCmd) Help_list() {
        println("Usage: list name")
    }

    func (this MyCmd) Do_list(name string) {
        println("name:", name, "received")
    }

    func (this MyCmd) Do_foo() {
    }

    c := cmd.New(new(MyCmd))
    c.Cmdloop()
*/
package cmd
