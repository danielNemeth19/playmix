# PlayMix

Project is aimed at creating randomized playlist files.

## Table of Contents
* [Usage](#usage)
    * [General Options](#general-options)
    * [Filtering Options](#filtering-options)
    * [Randomizing Options](#randomizing-options)
    * [Media Item Options](#media-item-options)
- [File format](#file-format)
- [Example XSPF format](#example-xspf-format)
- [VLC Extensions quick guide](#vlc-extensions-quick-guide)
- [TODO](#todo)

## Usage:
The tool expects an environment variable called `MEDIA_SOURCE` to be set as the `root` of all media files. 

This path will be recursively searched for all .mp4 files, respecting the below options.

### General Options
    -h, --help                  Print help and exists
    -ext                        If specified, collects unique file extensions
    -play                       If specified, playlist will be automatically played
    -fn                         Specifies the file name to use
                                (defaults to pl-test.xspf) 

### Filtering Options
    -fdate                      Only files after this date will be considered 
                                (defaults to "20000101")
    -tdate                      Only files up to this date will be considered
                                (defaults to "20300101")
    -mindur                     Consider only files with a duration longer
                                than the specified value (in seconds)
    -maxdur                     Consider only files with a duration shorter
    -include                    Folders to consider
    -skip                       Folders to skip
Both `include` and `skip` accepts a comma-separated list of folder names. The two option is mutually exclusive.


### Randomizing Options 
    -ratio                      Specifies the ratio of files to be included
                                (e.g. 80 means roughly 80%) 
    -stabilizer                 Specifies the interval at which elements are fixed
                                in place during shuffling (they still could be swapped)

### Media Item Options
    -options                    Allows for additional settings: it accepts a comma-
                                separated list of options
The available options are:
- `no-audio`: Excludes audio from the playlist.
- `start-time=<seconds>`: Sets the start time for the playlist in seconds.
- `stop-time=<seconds>`: Sets the stop time for the playlist in seconds.

Example usage:
```
-options="no-audio,start-time=30,stop-time=120"
```
This example will exclude audio, start the playlist at 30 seconds, and stop it at 120 seconds.

## File format
XSPF is a playlist in xml format - it is a free and open format.

## Example XSPF format
In its most simple form, a valid playlist looks like the below:

```xml
 <?xml version="1.0" encoding="UTF-8"?>
  <playlist version="1" xmlns="http://xspf.org/ns/0/">
     <trackList>
         <track><location>file:///mp3s/song_1.mp3</location></track>
         <track><location>file:///mp3s/song_2.mp3</location></track>
         <track><location>file:///mp3s/song_3.mp3</location></track>
     </trackList>
  </playlist>
```

## VLC Extensions quick guide
Source: [VideoLan Wiki](https://wiki.videolan.org/XSPF/)

XSFP supports extensions to allow applications to add special data. These extensions can appear in 'playlist' (node and item) or in 'track' (id and option).

Currently, extensions support the following elements:
* vlc:node
* vlc:item
* vlc:id
* vlc:option

The extensions vlc:node and vlc:item are used to specify how to display the playlist tree, which is not supported by standard XSPF.

### vlc:node

This element will be displayed as a node in the playlist. It appears as an extension of the **playlist** block (under playlist/extension). Only its name can be specified:
```xml
<vlc:node title="Node title">
  [list of vlc:item or vlc:node]
 </vlc:node>
```
### vlc:item

This element represents a playlist item (not a node). It appears as an extension of the **playlist** block (under playlist/extension). It contains only a track id (see below, vlc:id):
```xml
<vlc:item tid="42"/>
```
It seems vlc actually ignores the tid attribute and uses the track id from the vlc:id element to determine the order of the tracks in the playlist.
### vlc:id

This element specifies a track's id. It appears as an extension of the track **block** (under playlist/trackList/track/extension).
```xml
<vlc:id>42</vlc:id>
```
### vlc:option

This element allows you to add options to the input item. It appears as an extension of the **track** block (under playlist/trackList/track/extension).
```xml
<vlc:option>option-name</vlc:option>
```
Or, if the option has a value:

```xml
<vlc:option>option-name=option-value</vlc:option>
```
Example options:
* start-time (needs value)
* stop-time (needs value)
* no-audio

## TODO:
* supporting ratios per folder
* supporting timestamp filter for files (done)
* supporting regex patterns for file selection
* alternating files from folders (i.e. some kind of controlled randomization)
