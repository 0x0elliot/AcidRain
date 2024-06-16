"use client"

import { useEffect, useRef, useState } from "react"
import { siteConfig } from "@/app/siteConfig"
import cookies from 'nookies';

import { useEditor, EditorContent } from "@tiptap/react";
import { StarterKit } from "@syfxlin/tiptap-starter-kit";
import { Markdown } from 'tiptap-markdown';
import TextAlign from '@tiptap/extension-text-align'


import "./page.css"

export default function Post() {
  const [postInfo, setPostInfo] = useState({});
  const [postContent, setPostContent] = useState('');
  const [title, setTitle] = useState('What is the title?');

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
    content: postContent,
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
          console.log("Setting data", data);
          setPostInfo(data.post);
          console.log("Setting title as ", data.post.title);
          setTitle(data.post.title);
          setPostContent(data.post.content);
          if (editor !== null) {
            console.log("Setting editor content as ", data.post.content);
            editor.commands.setContent(data.post.content);
          }
        });
      } else if (response.status === 401) {
        // redirect to login
        window.location.href = '/login';
      }
    });
  }, []);

  useEffect(() => {
    if (editor !== null) {
      editor.commands.setContent(postContent);
    }
  }, [postContent]);

  return (
    
    <div>
      <h1>{title}</h1>

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
