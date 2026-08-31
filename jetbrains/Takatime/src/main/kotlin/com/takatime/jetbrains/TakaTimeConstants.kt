package com.takatime.jetbrains
import java.io.File


const val TAKATIME_VERSION = "v2.3.0"






object TakaTimeConfig {

    // Globally accessible function to grab the URI safely
    fun getMongoUri(): String? {
        val homeDir = System.getProperty("user.home")
        val configFile = File(homeDir, ".takatime.json")

        // If the file doesn't exist, they haven't set it up yet
        if (!configFile.exists()) return null

        return try {
            val content = configFile.readText()

            // Regex safely extracts the value without needing a heavy JSON parser
            val match = """"MONGO_URI"\s*:\s*"([^"]+)"""".toRegex().find(content)

            // Returns the matched URL, or null if it's empty/malformed
            match?.groupValues?.get(1)
        } catch (e: Exception) {
            println("TakaTime: Failed to read config - ${e.message}")
            null
        }
    }
}
