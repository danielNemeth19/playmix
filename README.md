# PlayMix

Project is aimed at creating randomized playlist files.

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
