use std::error::Error;
use std::fs;
use std::path::Path;

use rusqlite::{params, Connection};
use serde::ser::SerializeTuple;
use serde::{Serialize, Serializer};

type Result<T> = std::result::Result<T, Box<dyn Error>>;

#[derive(Serialize, Debug, Clone, Eq, PartialEq)]
#[serde(rename_all = "lowercase")]
enum DocType {
    Package,
    Func,
    Struct,
    Interface,
    Other,
}

impl DocType {
    fn from(ty: String) -> Self {
        match ty.as_str() {
            "package" => Self::Package,
            "func" => Self::Func,
            "struct" => Self::Struct,
            "interface" => Self::Interface,
            _ => Self::Other,
        }
    }
}

#[derive(Debug, Clone)]
struct DocItem {
    label: String,
    description: String,
    doc_type: DocType,
}

impl Serialize for DocItem {
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut ser = serializer.serialize_tuple(3)?;
        ser.serialize_element(&self.label)?;
        ser.serialize_element(&self.description)?;
        ser.serialize_element(&self.doc_type)?;
        ser.end()
    }
}

fn main() -> Result<()> {
    let mut data: Vec<DocItem> = vec![];
    let path = "gdocs.db";
    let db = Connection::open(path)?;
    let mut packages_stmt = db.prepare("select fullPath,simpleSpec from pkgs")?;
    for package in packages_stmt.query_map(params![], |row| {
        Ok(DocItem {
            doc_type: DocType::Package,
            label: row.get(0)?,
            description: row.get(1)?,
        })
    })? {
        let package = package?;
        let package_name = package.label.clone();
        let package_description = package.description.clone();

        let sql = "select label,simpleSpec,datatype from docs where pkg=(:package)";
        let mut docs_stmt = db.prepare(sql)?;
        for docs in docs_stmt.query_map(&[(":package", &package_name)], |row| {
            let doc_type = DocType::from(row.get(2)?);
            let label = if doc_type == DocType::Package {
                row.get(0)?
            } else {
                let l: String = row.get(0)?;
                format!("{}.{}", package_name, l)
            };
            Ok(DocItem {
                label,
                doc_type,
                description: row.get(1)?,
            })
        })? {
            let docs = docs?;
            data.push(docs);
        }
        data.push(DocItem {
            label: package_name,
            doc_type: DocType::Package,
            description: package_description,
        })
    }
    let contents = serde_json::to_string(&data)?;
    let path = Path::new("target/godocs.js");
    fs::write(path, format!("var searchIndex={};", contents))?;
    Ok(())
}
