# File Manager

[] - Open/ Create database file
[] - Read bytes from file
[] - Write bytes to file
[] - Flush changes to disk 
<!-- {I am thinking an exposed sync method instead of a hardcoded flush on every write} -->
[] - Get current file size
[] - Unit testing and file error handling

# Metadata

[] - Design metadata page
[] - Store DB version 
[] - Store page size
[] - Store root page id
[] - Store freelist page ID
<!-- {free list would be a package that would track which pages are empty in the database, thus not scaling the database endlessly} -->
[] - Load metadata on startup
<!-- {unsure how to handle that right now and in which package it'll be thus putting it here.} -->
[] - Update metadata on changes

# Page Manager

[] - Decide Page layout
<!-- {Could be a hardcoded value for now} -->
[] - Implement page structure
[] - Read page from disk
[] - Write page to disk
<!-- {would read or write on the bytes allocated to that page using file manager} -->
[] - Allocate New page

## Freelist
    [] - Free unused page
    [] - Reuse free pages

# Record Format

[] - Design on disk record layout 
<!-- {maybe something like --- 
 | key len | value len | key | value |} -->
[] - Encode records into byte
[] - Decode records from bytes
[] - Handle variable length keys and values

# B+ tree
:(

# WAL 

[] - Design WAL format
<!-- {something like ---
|operation| record format|
while we do have to make an exception for logging delete} -->
## Write commands to WAL
    [] - SET
    [] - DELETE
    [] - UPDATE

[] - Replay WAL during startup
[] - Recover db on crash
[] - Clear old WAL after a fixed checkpoint
<!-- We can also play with rollback and truncation here -->

# Storage ENgine
## Implement
    [] - PUT
    [] - GET
    [] - DELEte
    [] - UPDATE

and more... :]

# MVCC
[] - Design revision format
[] - Add revision number to every write {might require changing the record or WAL, completely clueless right now}
[] - Store multiple versions of a key
[] - Read latest version
[] - read specific old version
[] - delete old revisions after certain commits or logging

# RAFT

[] - DEsign raft logs
[] - Implement and test states (follower, candidate, leader)
[] - Implement election
[] - Implement heartbeat
[] - Implement Request vote
[] - Replicate log entries
[] - majority voting to commit entries
[] - committed entries to be synced to storage through the exposed sync function
[] - Recover node after restart

# gRPC
<!-- already done needs refinement, maybe :) -->
