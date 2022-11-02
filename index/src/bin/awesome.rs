use std::borrow::Cow;
use std::collections::HashMap;
use std::fs;
use std::path::Path;

use regex::Regex;
use serde::ser::SerializeTuple;
use serde::{Serialize, Serializer};

const AWESOME_MARKDOWN_URL: &str =
    "https://raw.githubusercontent.com/avelino/awesome-go/master/README.md";

type Result<T> = std::result::Result<T, Box<dyn std::error::Error>>;

#[derive(Debug)]
struct Awesome<'a> {
    name: &'a str,
    url: &'a str,
    description: Cow<'a, str>,
    category: Option<&'a str>,
}

impl Serialize for Awesome<'_> {
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut ser = serializer.serialize_tuple(4)?;
        // ser.serialize_element(&self.get_package_domain())?;
        ser.serialize_element(&self.name)?;
        ser.serialize_element(&self.url)?;
        ser.serialize_element(truncate(self.description.as_ref(), 100))?;
        if let Some(category) = self.category {
            ser.serialize_element(category)?;
        }
        ser.end()
    }
}

fn truncate(s: &str, max_chars: usize) -> &str {
    match s.char_indices().nth(max_chars) {
        None => s,
        Some((idx, _)) => &s[..idx],
    }
}

fn assert_link(text: &str) -> bool {
    text.contains(|c: char| ['[', ']'].contains(&c))
}

#[tokio::main]
async fn main() -> Result<()> {
    let mut data = vec![];
    let mut others = HashMap::new();
    let r = Regex::new(r"\[(?P<name>.*)]\(https?.*\)").unwrap();

    let regex: Regex =
        Regex::new(r"^\*\s*\[(?P<name>[^\[\]]*)]\((?P<url>https?\S*)\)\s*-?\s*(?P<desc>.*)?$")
            .unwrap();
    let markdown = reqwest::get(AWESOME_MARKDOWN_URL).await?.text().await?;
    // let markdown = include_str!("file.md");
    let mut current_category = "";
    for line in markdown.split_terminator('\n').collect::<Vec<&str>>() {
        if line.starts_with("## ") {
            current_category = (line[2..]).trim();
        } else if let Some(cap) = regex.captures(line.trim()) {
            let name = cap.name("name").unwrap().as_str();
            assert!(!assert_link(name));

            let url = cap.name("url").unwrap().as_str();
            assert!(!assert_link(url));

            let description = r.replace_all(cap.name("desc").unwrap().as_str(), "$name");
            match current_category {
                "Conferences" | "E-Books" | "Meetups" | "Social Media" => {
                    let category = match current_category {
                        "Conferences" => "conf",
                        "E-Books" => "book",
                        "Meetups" => "meetup",
                        "Social Media" => "social",
                        _ => "",
                    };
                    let vec = others.entry(category).or_insert(vec![]);
                    vec.push(Awesome {
                        name,
                        url,
                        description,
                        category: None,
                    });
                }
                _ => {
                    data.push(Awesome {
                        name,
                        url,
                        description,
                        category: Some(current_category),
                    });
                }
            }
        }
    }
    let contents = serde_json::to_string(&data)?;
    let path = Path::new("target/awesome.js");
    fs::write(path, format!("var awesomeIndex={};", contents))?;

    let contents = serde_json::to_string(&others)?;
    let path = Path::new("target/others.js");
    fs::write(path, format!("var othersIndex={};", contents))?;
    Ok(())
}
