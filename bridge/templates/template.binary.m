% Flywheel
classdef Flywheel
    % Flywheel class enables user to communicate with Flywheel platform
    properties
        key     % key - API Key assigned through the Flywheel UI
        folder  % folder - folder where the SDK code is located
    end
    methods
        function obj = Flywheel(apiKey)
            % Usage Flywheel(apiKey)
            %  apiKey - API Key assigned for each user through the Flywheel UI
            %          apiKey must be in format <domain>:<API token>
            C = strsplit(apiKey, ':');
            % Check if key is valid
            if length(C) < 2
                ME = MException('FlywheelException:Invalid', 'Invalid API Key');
                throw(ME)
            end
            obj.key = apiKey;
            % Check if JSONlab is in path
            if ~exist('savejson')
                ME = MException('FlywheelException:JSONlab', 'JSONlab function savejson is not loaded. Please install JSONlab and add to path.')
                throw(ME)
            end
            if ~exist('loadjson')
                ME = MException('FlywheelException:JSONlab', 'JSONlab function loadjson is not loaded. Please install JSONlab and add to path.')
                throw(ME)
            end

            [folder, name, ext] = fileparts(mfilename('fullpath'));
            obj.folder = folder;

            %%% TODO: use this code below to determine which .so and .h to load
            %if ismac
            % Code to run on Mac plaform
            %elseif isunix
            % Code to run on Linux plaform
            %elseif ispc
            % Code to run on Windows platform
            %else
            %    disp('Platform not supported')
            %end
            % Suppress Max Length Warning
            warningid = 'MATLAB:namelengthmaxexceeded';
            warning('off',warningid);
        end

        % TestBridge
        function cmdout = testBridge(obj, s)
            [status,cmdout] = system([obj.folder '/sdk TestBridge ' s]);
        end

        %
        % AUTO GENERATED CODE FOLLOWS
        %

        {{range .Signatures}}% {{camel2lowercamel .Name}}
        function result = {{camel2lowercamel .Name}}(obj{{range .Params}}, {{.Name}}{{end}})
            statusPtr = libpointer('int32Ptr',-100);
            {{if ne .ParamDataName ""}}oldField = 'id';
            newField = 'x0x5F_id';
            {{.ParamDataName}} = Flywheel.replaceField({{.ParamDataName}},oldField,newField);
            {{.ParamDataName}} = savejson('',{{.ParamDataName}});
            {{end -}}
            [status,cmdout] = system([obj.folder '/sdk {{.Name}} ' obj.key ' ' {{range .Params}} '"' {{.Name}} '" '{{end -}}]);

            result = Flywheel.handleJson(status,cmdout);
        end
        {{end}}
        % AUTO GENERATED CODE ENDS
    end
    methods (Static)
        function version = getSdkVersion()
            version = '{{.Version}}';
        end
        function structFromJson = handleJson(statusPtr,ptrValue)
            % Handle JSON using JSONlab
            statusValue = statusPtr;

            % If status indicates success, load JSON
            if statusValue == 0
                % Interpret JSON string blob as a struct object
                loadedJson = loadjson(ptrValue);
                % loadedJson contains status, message and data, only return
                %   the data information.
                dataFromJson = loadedJson.data;
                %  Call replaceField on loadedJson to replace x0x5F_id with id
                structFromJson = Flywheel.replaceField(dataFromJson,'x0x5F_id','id');
            % Otherwise, nonzero statusCode indicates an error
            else
                % Try to load message from the JSON
                try
                    loadedJson = loadjson(ptrValue);
                    msg = loadedJson.message;
                    ME = MException('FlywheelException:handleJson', msg);
                % If unable to load message, throw an 'unknown' error
                catch ME
                    msg = sprintf('Unknown error (status %d).',statusValue);
                    causeException = MException('FlywheelException:handleJson', msg);
                    ME = addCause(ME,causeException);
                    rethrow(ME)
                end
                throw(ME)
            end
        end
        function newStruct = replaceField(oldStruct,oldField,newField)
            % Replace a field within a struct or a cell array of structs
            % Check if variable is a cell
            if iscell(oldStruct)
                % Initialize newStruct as a copy of the oldStruct
                newStruct = oldStruct;
                for k=1:length(oldStruct)
                    f = fieldnames(oldStruct{k});
                    % Check if oldField is a fieldname in oldStruct
                    if any(ismember(f, oldField))
                        [oldStruct{k}.(newField)] = oldStruct{k}.(oldField);
                        newStruct{k} = rmfield(oldStruct{k},oldField);
                    else
                        newStruct{k} = oldStruct{k};
                    end
                end
            % Check if variable is a struct
            elseif isstruct(oldStruct)
                % Replace a fieldname within a struct object
                f = fieldnames(oldStruct);
                % Check if oldField is a fieldname in oldStruct
                if any(ismember(f, oldField))
                    [oldStruct.(newField)] = oldStruct.(oldField);
                    newStruct = rmfield(oldStruct,oldField);
                else
                    newStruct = oldStruct;
                end
            % If not, newStruct is equal to oldStruct
            else
                newStruct = oldStruct;
            end
        end
    end
end
