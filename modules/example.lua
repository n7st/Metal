-- A basic example script which handles an incoming command ("hello") and
-- returns string output to be sent back as a message.

-- Options are the module's options. These are parsed by the bot to understand
-- when the script should be run, and what should be targeted.
--
-- Setting ["message"] to true requires a "Process" function, and setting
-- ["command"] to true requires a "Command" function.
function Options()
    return {
        ['message']  = true,
        ['command']  = true,
        ['commands'] = {
            -- The commands table is a map of command to method to call. This
            -- means that adding ['!translate_goodbye'] = 'Translate_Goodbye'
            -- would attempt to run Translate_Goodbye() if a user sent a message
            -- starting with "!translate_goodbye".
            ['!translate_hello'] = 'Translate_Hello',
        },
    }
end

-- The Process function should process a full input message. This could be used
-- for simple string processing (as seen here), or for something that involves
-- parsing the input, like a spam detection script.
--
-- The function's return value will be sent back to where the message was
-- received.
function Process(message)
    if message == 'hello' then
        return 'Hello, world'
    elseif message == 'goodbye' then
        return 'Goodbye, world'
    end

    -- Providing a default "return" would cause the bot to reply to every
    -- message (not recommended).
    --
    -- return 'Spammy!'
end

-- Command methods must be global.
function Translate_Hello(to_language)
    if to_language == '' then
        return
    end

    string.lower(to_language)

    if to_language == 'french' then
        return 'Bonjour'
    elseif to_language == 'german' then
        return 'Guten tag'
    end

    return "I don't know that language."
end