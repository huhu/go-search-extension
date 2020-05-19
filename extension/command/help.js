class HelpCommand extends Command {
    constructor() {
        super("help", "Show the help messages.");
    }

    onExecute(arg) {
        const value = ([
            `Prefix ${c.match(":")} to execute command (:help, :history, etc)`,
            `Prefix ${c.match("!")} to search packages exclusively`,
        ]);
        return value.map((description, index) => {
            return {content: `${index + 1}`, description};
        });
    }
}