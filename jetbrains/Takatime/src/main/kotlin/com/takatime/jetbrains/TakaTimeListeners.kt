package com.takatime.jetbrains

import com.intellij.openapi.editor.Document
import com.intellij.openapi.editor.EditorFactory
import com.intellij.openapi.editor.event.DocumentEvent
import com.intellij.openapi.editor.event.DocumentListener
import com.intellij.openapi.editor.event.EditorFactoryEvent
import com.intellij.openapi.editor.event.EditorFactoryListener
import com.intellij.openapi.fileEditor.FileDocumentManagerListener
import java.util.WeakHashMap

// 1. Save Listener
class TakaTimeSaveListener : FileDocumentManagerListener {
    override fun beforeDocumentSaving(document: Document) {
        TakaTimeHeartbeat.handleHeartbeat(document)
    }
}

// 2. Typing Listener
class TakaTimeEditorFactoryListener : EditorFactoryListener {

    private val typingListener = object : DocumentListener {
        override fun documentChanged(e: DocumentEvent) {
            TakaTimeHeartbeat.handleHeartbeat(e.document)
        }
    }

    override fun editorCreated(event: EditorFactoryEvent) {
        val editor = event.editor
        TakaTimeBinaryManager.checkAndDownloadIfNeeded(editor.project)

        // JetBrains 2026.2+ compatibility: use a disposable-based registration
        // so the listener is removed when the project/editor is disposed.
        editor.document.addDocumentListener(typingListener, editor.disposable)

    }

    override fun editorReleased(event: EditorFactoryEvent) {
        // Listener lifecycle is now managed by the disposable, so this is no longer needed.
    }
}