class HelpCommand extends Command {
    constructor() {
        super("help", "Show the help messages.");
    }

    onExecute(arg) {
        const value = [
            `Prefix ${c.match(":")} to execute command (:help, :history, etc)`,
            `Prefix ${c.match("!")} to search packages exclusively, prefix ${c.match("!!")} to open the repository`,
            `Prefix ${c.match("$")} to search a curated list of awesome Go frameworks, libraries and software`,
        ];
        return this.wrap(value);
    }
}