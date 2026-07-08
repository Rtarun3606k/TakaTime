package com.takatime.jetbrains

import com.intellij.openapi.editor.Document
import com.intellij.openapi.fileEditor.FileDocumentManager
import com.intellij.openapi.project.ProjectLocator
import com.jetbrains.rd.generator.nova.PredefinedType
import java.nio.file.Paths
import com.intellij.openapi.application.ApplicationInfo

object TakaTimeUploader {

    fun spawnProcess(document: Document, duration: String): Boolean {
        // 1. Check Config
        val mongoUri = TakaTimeConfig.getMongoUri()
        if (mongoUri.isNullOrBlank()) return false

        // 2. Extract the file
        val file = FileDocumentManager.getInstance().getFile(document)
        if (file == null || !file.isInLocalFileSystem) return false
        if (file.path.contains("/.git/")) return false

        // 3. Extract the Telemetry Data
        val filePath = file.path

        val project = ProjectLocator.getInstance().guessProjectForFile(file)
        val projectName = project?.basePath?.let { java.io.File(it).name }
            ?: project?.name
            ?: "UnknownProject"

        // 4. Safely build the binary path
        val isWindows = System.getProperty("os.name").lowercase().contains("win")
        val ext = if (isWindows) ".exe" else ""

        val takaUploadPath = Paths.get(
            System.getProperty("user.home"),
            ".takatime",
            "bin",
            "taka-upload-$TAKATIME_VERSION$ext"
        ).toString()
        val applicationName = ApplicationInfo.getInstance().versionName ?: "JetBrains"
        // 5. Fire the Go CLI!
        return try {
            ProcessBuilder(
                takaUploadPath,
                "-file", filePath,
                "-project", projectName,
                "-uri", mongoUri,
                "-duration", duration,
                "-editor", applicationName,
            ).start()

            println("TakaTime: Heartbeat Sent -> Project: $projectName | File: $filePath | uri: $takaUploadPath | duration: $duration | application: $applicationName")
            true // Success! Timer will reset.
        } catch (e: Exception) {
            println("TakaTime: Failed to execute CLI - ${e.message}")
            false // Failed! Timer will NOT reset.
        }
    }
}
