use std::collections::HashMap;
use std::error::Error;
use std::fs;
use std::path::Path;

use rusqlite::{params, Connection};
use serde::ser::SerializeTuple;
use serde::{Serialize, Serializer};

type Result<T> = std::result::Result<T, Box<dyn Error>>;

#[derive(Debug, Clone)]
struct Package {
    full_path: String,
    description: String,
    version: String,
}

impl Serialize for Package {
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut ser = serializer.serialize_tuple(2)?;
        // ser.serialize_element(&self.get_package_domain())?;
        ser.serialize_element(&self.description)?;
        ser.serialize_element(&self.version)?;
        ser.end()
    }
}

fn main() -> Result<()> {
    let mut data: HashMap<String, Package> = HashMap::new();
    let path = "pkgs.db";
    let db = Connection::open(path)?;
    let statement = "select fullPath,simpleSpec,version from pkgs \
                            where libtype in (\"package\", \"directory\", \"command\") \
                            and version not like 'v0.0.0%' \
                            order by uid limit 20000";
    let mut packages_stmt = db.prepare(statement)?;
    for package in packages_stmt.query_map(params![], |row| {
        Ok(Package {
            full_path: row.get(0)?,
            description: row.get(1)?,
            version: row.get(2)?,
        })
    })? {
        let package = package?;
        data.insert(package.full_path.to_string(), package);
    }
    println!("Total package is: {}", data.len());
    let contents = serde_json::to_string(&data)?;
    let path = Path::new("target/packages.js");
    fs::write(path, format!("var pkgs={};", contents))?;
    Ok(())
}
