"use client"

import { useEffect, useRef } from "react"
import { siteConfig } from "@/app/siteConfig"
import cookies from 'nookies';

import { useEditor, EditorContent } from "@tiptap/react";
import { StarterKit } from "@syfxlin/tiptap-starter-kit";
import { Markdown } from 'tiptap-markdown';
import TextAlign from '@tiptap/extension-text-align'


import "./page.css"

export default function Post() {

  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        // disable
        emoji: false,
        // configure
        heading: {
          levels: [1, 2],
        },
      }),
      Markdown,
      TextAlign.configure({
        types: ['heading', 'paragraph'],
      }),
    ],
  });

  if (editor !== null) {
    const markdownOutput = editor.storage.markdown.getMarkdown();
  }

  useEffect(() => {
    let accessToken = cookies.get(null).access_token;

    const urlParams = new URLSearchParams(window.location.search);
    const postId = urlParams.get('post_id');

    fetch(`${siteConfig.baseApiUrl}/api/post/private/${postId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': `Bearer ${accessToken}`,
      },
    })
    .then((response) => {
      if (response.ok) {
        response.json().then((data) => {
          console.log(data);
        });
      } else if (response.status === 401) {
        // redirect to login
        window.location.href = '/login';
      }
    });
  }, []);

  return (
    
    <div>
      <h1>Write Post</h1>

      <div className="control-group">
        <div className="button-group">
          <button
            onClick={() => editor?.chain().focus().setTextAlign('left').run()}
            className={editor?.isActive({ textAlign: 'left' }) ? 'is-active' : ''}
          >
            Left
          </button>
          <button
            onClick={() => editor?.chain().focus().setTextAlign('center').run()}
            className={editor?.isActive({ textAlign: 'center' }) ? 'is-active' : ''}
          >
            Center
          </button>
          <button
            onClick={() => editor?.chain().focus().setTextAlign('right').run()}
            className={editor?.isActive({ textAlign: 'right' }) ? 'is-active' : ''}
          >
            Right
          </button>
          <button
            onClick={() => editor?.chain().focus().setTextAlign('justify').run()}
            className={editor?.isActive({ textAlign: 'justify' }) ? 'is-active' : ''}
          >
            Justify
          </button>
          <button onClick={() => editor?.chain().focus().unsetTextAlign().run()}>Unset text align</button>
        </div>
      </div>
      
      <div className="editor-container" style={{ marginBottom: '1rem' }}>
        <EditorContent editor={editor} className="editor-content" />
      </div>
      
    </div>
  )
}
