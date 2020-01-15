# -*- coding: utf-8 -*-

import os
import sys
import json
import util

# The attribute mapping for docsearch.
#
# The 'tags' mapping's value will be a list, so values from multiple
# taxonomies can be put into the docsearch `tags` property.
docsearch_mapping = { "content": "content",
                      "url": "url",
                      "tags": ["tags", "categories"]
                    }

# default values for the weight property
docsearch_weight = { "position": 1,
                     "level": 10,
                     "page_rank": 0
                    }

# top hierarchy level. Sections will be used for additional
# levels
base_level = "Hugo Site"

def create_index_list(walk_dir, bsURL):
    """Create a list of index entries starting from the directory walk_dir"""
    base_url = bsURL
    global base_level, docsearch_mapping, docsearch_weight

    # used to store all indexed item (markdown files)
    index_list = []
    
    # give items an objectID so that the index file can be uploaded
    # to algolia again, overwriting previous index items
    objectID = 0

    for root, subdirs, files in os.walk(walk_dir):
        for filename in files:

            # index md files
            if filename.endswith(".md"):
                objectID += 1
                filepath = os.path.join(root, filename)

                subpath = root[len(walk_dir):].rstrip(os.sep)
                subpaths = subpath.lstrip(os.sep).split(os.sep)
                
                # index.md have special URLs
                if filename != "index.md":
                    subpaths.append(filename[:-3])
                
                # set up list for the hierarchy of the markdown file
                hierarchy_list = [base_level]
                hierarchy_list.extend(subpaths)
                
                # construct the url of the markdown file
                url_subpath = "/".join(subpaths)
                url = base_url + "/" + url_subpath + "/"

                sys.stderr.write("Indexing '" + filepath + "' (" + url + "\n")

                # get data from the file (frontmatter and content)
                filedata = parse_md(filepath)

                # create index entry
                indexed_item = {'objectID': objectID, 'url': url.lower() }

                # map filedata to docsearch structure
                for docsearch_key, filedata_key in docsearch_mapping.items():
                    
                    # plain mappings, configured at the top of the script
                    if type(filedata_key) == str and filedata_key in filedata.keys():
                        indexed_item[docsearch_key] = filedata[filedata_key]

                    # if the mapping value is a list, assume that the frontmatter data
                    # of the keys in the list are also lists. Combine the lists values of each
                    # frontmatter property into a list and set the docsearch property as
                    # to this combined value list (used for the "tags" property. see
                    # comment in the beginning of the script
                    elif type(filedata_key) == list:
                        aggregated = []
                        for filedata_subkey in filedata_key:
                            if filedata_subkey in filedata.keys():
                                aggregated.extend(filedata[filedata_subkey])
                        indexed_item[docsearch_key] = aggregated

                    # hierarchy and hierarchy_complete properties
                    hierarchy = create_empty_hierarchy()
                    hierarchy_complete = create_empty_hierarchy()
                    for level in range(7):
                        if level < len(hierarchy_list):
                            hierarchy["lvl" + str(level)] = hierarchy_list[level]
                            hierarchy_complete["lvl" + str(level)] = " > ".join(hierarchy_list[:level])
                    indexed_item["hierarchy"] = hierarchy
                    indexed_item["hierarchy_complete"] = hierarchy_complete

                    # hierarchy_radio and type
                    hierarchy_radio = create_empty_hierarchy()
                    max_lvl = len(subpaths) - 1
                    hierarchy_radio["lvl" + str(max_lvl)] = subpaths[max_lvl]
                    indexed_item["hierarchy_radio"] = hierarchy_radio
                    indexed_item["type"] = "lvl" + str(max_lvl)

                    # anchor and weight. anchors are not considered
                    indexed_item["anchor"] = None
                    indexed_item["weight"] = docsearch_weight

                index_list.append(indexed_item)
    sys.stderr.write("Done indexing .md files in '" + walk_dir + "'" + "\n")
    return index_list

def create_empty_hierarchy():
    """Create a empty hierarchy structure (dict)."""
    empty_hierarchy = {}
    for level_index in range(7):
        empty_hierarchy["lvl" + str(level_index)] = None
    return empty_hierarchy


if __name__ == '__main__':
    if len(sys.argv) != 2:
        sys.stderr.write("ERROR: Please give me your config file!")
        sys.exit(1)
    
    cfgFileName = sys.argv[1]
    fileExt = cfgFileName.split(".")[1]
    z = infoGet(fileExt)
    
    # gather index data
    index_list = create_index_list("./"+z["contentdir"], z["baseURL"])

    # output the index as readable json to stdout. does not escape UTF-8 characters
    #sys.stdout.write(json.dumps(index_list, ensure_ascii=False, indent=2))
    with open("./contIndex.json", "w") as f:
        f.write(json.dumps(index_list, ensure_ascii=False, indent=2))
