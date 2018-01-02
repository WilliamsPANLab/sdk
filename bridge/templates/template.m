% Flywheel
classdef Flywheel
    % Flywheel class enables user to communicate with Flywheel platform
    properties
        key     % key - API Key assigned through the Flywheel UI
        root    % root - Whether or not this client is in management mode
    end
    methods
        function obj = Flywheel(apiKey, root)
            % Usage Flywheel(apiKey)
            %  apiKey - API Key assigned for each user through the Flywheel UI
            %          apiKey must be in format <domain>:<API token>
            %  root - Set to 'true' to indicate that the client should run in manage mode
            C = strsplit(apiKey, ':');
            % Check if key is valid
            if length(C) < 2
                ME = MException('FlywheelException:Invalid', 'Invalid API Key');
                throw(ME)
            end
            obj.key = apiKey;

            % Set root mode
            if exist('root', 'var') && root
                obj.root = 1;
            else
                obj.root = 0;
            end

            % Check if JSONio is in path
            if ~exist('jsonread')
                ME = MException('FlywheelException:JSONio', 'JSONio function jsonread is not loaded. Please install JSONio and add to path.')
                throw(ME)
            end
            if ~exist('jsonwrite')
                ME = MException('FlywheelException:JSONio', 'JSONio function jsonwrite is not loaded. Please install JSONio and add to path.')
                throw(ME)
            end
            % Load flywheel shared library
            Flywheel.loadBridge()
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

        function ptrValue = testBridge(obj, s)
            Flywheel.loadBridge()
            % Call bridge
            ptrValue = calllib('flywheelBridge','TestBridge',s);
        end

        %
        % AUTO GENERATED CODE FOLLOWS
        %

        {{range .Signatures}}{{ $length := .LastParamIndex }}
        function result = {{camel2lowercamel .Name}}(obj{{range .Params}}, {{.Name}}{{end}})
            % {{camel2lowercamel .Name}}({{range $ind, $val := .Params}}{{.Name}}{{if lt $ind $length}}, {{end}}{{end -}})

            statusPtr = libpointer('int32Ptr',-100);
            {{if ne .ParamDataName ""}}oldField = 'id';
            newField = 'x0x5Fid';
            {{.ParamDataName}} = Flywheel.replaceField({{.ParamDataName}},oldField,newField);
            % Indicate to JSONio to replace hex values with corresponding character, i.e. 'x0x5F' -> '_' and '0x2D' -> '-'
            opts = struct('replacementStyle','hex');
            {{.ParamDataName}} = jsonwrite({{.ParamDataName}},opts);
            {{end -}}
            pointer = calllib('flywheelBridge','{{.Name}}',obj.key,obj.root,{{range .Params}}{{.Name}},{{end -}} statusPtr);
            result = Flywheel.handleJson(statusPtr,pointer);
        end
        {{end}}
        % AUTO GENERATED CODE ENDS
    end
    methods (Static)
        function version = getSdkVersion()
            version = '{{.Version}}';
        end
        function structFromJson = handleJson(statusPtr,ptrValue)
            % Handle JSON using JSONio
            statusValue = statusPtr.Value;
            % If status indicates success, load JSON
            if statusValue == 0
                % Interpret JSON string blob as a struct object
                loadedJson = jsonread(ptrValue);
                % loadedJson contains status, message and data, only return
                %   the data information.
                dataFromJson = loadedJson.data;
                %  Call replaceField on loadedJson to replace x_id with id
                structFromJson = Flywheel.replaceField(dataFromJson,'x_id','id');
            % Otherwise, nonzero statusCode indicates an error
            else
                % Try to load message from the JSON
                try
                    loadedJson = jsonread(ptrValue);
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
        % loadBridge
        function [] = loadBridge()
            % Load flywheel shared library
            if not(libisloaded('flywheelBridge'))
                % loading in flywheelBridge.so file
                try
                    loadlibrary('flywheelBridge', 'flywheelBridgeSimple.h')
                catch ME
                    % Add some user-friendly message for OSX
                    if (ismac && strcmp(ME.identifier, 'MATLAB:mex:NoCompilerFound_link'))
                        % TODO: improve message to the user
                        causeException = MException('FlywheelException:NoCompiler', 'Please install Xcode.');
                        ME = addCause(ME, causeException);
                    end
                    rethrow(ME)
                end
            end
        end
    end
end
