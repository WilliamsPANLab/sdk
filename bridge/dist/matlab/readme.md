# Flywheel Matlab SDK
Software Development Kit for Matlab

Use the Flywheel Matlab SDK to do the following:
* download and upload files to the Flywheel platform
* start and monitor Flywheel Gear runs
* create groups, projects, sessions, and acquisitions
* create tags and notes for groups, projects, sessions, and acquisitions

### Requirements
Requires [JSONio](https://github.com/gllmflndn/JSONio) Matlab package.

### Convention Note
Due to Matlab restrictions, structure member names cannot begin with an underscore. 
Therefore the `_id` field is mapped to `id` automatically by the SDK.

### Example Usage
Before starting, add JSONio and Flywheel Matlab SDK to project path.

```matlab
addpath('/path/to/JSONio/');
addpath('/path/to/FlywheelSDK');
```

Create a Flywheel object using the API key (located on User Profile on Flywheel UI).
```matlab
fw = Flywheel('key');
```

#### downloading a file from Flywheel

Get all Projects from the Flywheel Platform.
```matlab
projects = fw.getAllProjects();
```

Select id of desired project from `projects` cell.
```matlab
projectId = projects{1}.id;
```

Get all Sessions within project.
```matlab
projectSessions = fw.getProjectSessions(projectId);
```

Select id of desired session from `projectSessions` cell.
```matlab
sessionId = projectSessions{1}.id;
```

Get all Acquisitions within session.
```matlab
sessionAcqs = fw.getSessionAcquisitions(sessionId);
```

Select id of desired acquisition from sessionAcqs cell, select the name of the file desired to be downloaded, and define pathname on local machine to download the file.
```matlab
acqId = sessionAcqs{1}.id;
name = sessionAcqs{1}.files{1}.name;
path = '/path/of/filename';
fw.downloadFileFromAcquisition(acqId, name, path);
```

The file should now be downloaded to `path`.


#### uploading a file from Flywheel

Define the id of the acquisition to upload a file to and the location of the file on local machine. Use the same method to define the acquisition id as used above to download a file.
```matlab
acqId = 'acquisitionid'
newFile = '/path/to/file/to/upload';
fw.uploadFileToAcquisition(acqId, newFile);
```

The file will now be on the Flywheel platform within the designated acquisition.


#### adding notes and tags
Define the session id to add a note and tag to. Use the same method above to determine desired session id. NOTE: Tags and notes can be assigned to a project, session, or an acquisition.
```matlab
sessionId = 'sessionid'
fw.addSessionTag(sessionId, 'tag');
fw.addSessionNote(sessionId, 'This is a note');
```
