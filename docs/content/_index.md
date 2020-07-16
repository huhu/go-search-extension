+++
title = "Go Search Extension"
sort_by = "weight"
+++


Hi all. I'm so honored to introduce my product: **Go Search Extension,** a handy browser plugin to helps every Go developer search docs and package in the address bar instantly. 

So how does it works? Easy! Just input keyword **go** + **Space** in the address bar to activate the extension, then input any keyword you wanna search, the extension will respond to the related result in an instant! It's hugely fast just like a **millisecond-level** search! Why it so fast? Because we build the whole standard docs and packages as an offline index file.

> I also the creator of **Rust Search Extension**, which gets a lot of traction from the rust community. And more programming version (such as C/C++, Javascript) is just in the plan. See [https://github.com/huhu](https://github.com/huhu) if you want to know more.

Let's explain it in more detail:

### Search std docs

The whole standard library is searchable. No matter a `package`, `func`, or `interface`, the result you desired is just a keystroke away. Select one and enter, you'll be redirected to the proper [pkg.go.dev](http://pkg.go.dev) page effortlessly.

![](/std.png)

### Search top 20K packages

Just searching the std library is not enough. So we crawled top starred Go packages from Github and built the index file for TOP 20K, then everyone can search those great packages easily and instantly. 

![](/package.png)

Hope you like it, we bring several prefix sign to help you search the corresponding kind of content exclusively. For example, the **!** is package searching prefix, the **!!** is repository mode prefix, the **$** is the awesome golang list searching prefix (See the following part). If no prefix sign is specified, the default search is the standard library. :)

### Quick jump to git repository

While we searching the top package, we also hope we can jump to the package repository directly. Yeah, this is possible! Just prefix **!!** before the keyword, you'll open the repository page instead of the default [pkg.go.dev](http://pkg.go.dev) page. However, we only support Github, Gitlab, Bitbucket, Gitea, etc.

![](/repository-mode.png)

### Awesome golang list search

[Awesome Golang](https://github.com/avelino/awesome-go) is an aweomse resource itself for searching. Hence, we indexed the whole list! 

Prefix **$** before the keyword, you'll get the desired result from the awesome list.

![](/awesome.png)

### Builtin commands

The command system brings a handy set of useful and convenient commands to you. Each command starts with a **:** (colon), followed by the name, and function differently in individual. Those commands including but not limited to:

- `:help` - Show the help messages.
- `:book` - Show all Golang e-books.
- `:conf` - Show Golang conferences.
- `:meetup` - Show Golang meetups.
- `:social` - Show Golang social medias.
- `:history` - Show your local search history

![](/commands.png)

## **Page down/up easily**

You can press `space` after the keyword, then increase or decrease the number of **-** (hyphen) to page down or page up.
