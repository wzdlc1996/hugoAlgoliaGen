# -*- coding: utf-8 -*-

import toml
import yaml
import stopwords
import sys
import re

noNeedPatts = [
            r"`\$.*?\$`",
            r"<div>.*?</div>",
            r"<br>",
            r"<em>.*?</em>",   
            r"\$\$",
            r"\n",
            r"\\",
            r"<div>.*?</div>",
            r"#"
        ]

def infoGet(method = "toml"):
    """
    Get the basic information from the config file in hugo-site root dir.
    Input:
        method: file extension name of the config file. 
                should be one of the set:
                ["toml", "yaml", "json"]
    Output:
        z: the dict format data of the config file.
    """
    if (method == "toml"):
        with open("./config.toml", "r") as f:
            z = toml.load(f)
    else:
        print("Other formats are not supported yet!")
        raise ValueError("Give me toml pls")
    
    return z

def mdParser(filepath):
    """Parse a markdown file. Reads yaml front matter."""
    yaml_string = ""
    in_yaml = None
    content = ""
    with open(filepath, "r") as datafile:
        for line in datafile:
            if line.startswith("---"):
                if in_yaml:
                    in_yaml = False
                else:
                    in_yaml = True
                    continue
            elif in_yaml == True:
                yaml_string += line
            else:
                content += line
    md_data = yaml.load(yaml_string)
    if not "content" in md_data.keys():
        for patt in noNeedPatts:
            content = re.sub(patt, " ", content)
        wordlis = [wd for wd in re.findall(r"[A-Za-z]+", content) if wd not in stopwords.stopwords]
        #content = content.split(r"/s+")
        md_data["content"] = " ".join(wordlis)
    else:
        sys.stderr.write("ERROR: Could not store content for '" + filepath + "'\n")
    return md_data