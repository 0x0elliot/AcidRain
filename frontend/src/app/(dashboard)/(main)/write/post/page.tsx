"use client"

import { useEffect, useState } from "react"
import { siteConfig } from "@/app/siteConfig"
import cookies from 'nookies';

import { useEditor, EditorContent } from "@tiptap/react";
import { StarterKit } from "@syfxlin/tiptap-starter-kit";
import { Markdown } from 'tiptap-markdown';
import TextAlign from '@tiptap/extension-text-align'
import ImageResize from 'tiptap-extension-resize-image';

import toast, { Toaster } from 'react-hot-toast';


import "./page.css"

export default function Post() {
  const [postInfo, setPostInfo] = useState({});
  const [isEditing, setIsEditing] = useState(false);
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
      ImageResize,
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
          setPostInfo(data.post);
          setTitle(data.post.title);
          setPostContent(data.post.content);
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


  const handleEditClick = () => {
    setIsEditing(true);
  };

  const handleInputChange = (e) => {
    setTitle(e.target.value);
  };

  const handleBlur = () => {
    setIsEditing(false);
  };

  const handleSave = () => {
    // get current content
    if (editor !== null) {
      let content = editor.getHTML();

      let accessToken = cookies.get(null).access_token;

      fetch(`${siteConfig.baseApiUrl}/api/post/private/set`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'Authorization': `Bearer ${accessToken}`,
        },
        body: JSON.stringify({
          id: postInfo.id,
          title: title,
          content: content,
        }),
      }).then((response) => {
        if (response.status === 200) {
          toast.success('Post saved successfully');
        } else {
          toast.error('Failed to save post');
        }
      }
      );

    }
  };

  return (
    <div>
      <div style={{ display: 'grid', gridTemplateColumns: '1fr auto' }}>
        <Toaster />
        <div>
          {isEditing ? (
            <input
              type="text"
              value={title}
              onChange={handleInputChange}
              onBlur={handleBlur}
              autoFocus
            />
          ) : (
            <h1 onClick={handleEditClick}>{title}</h1>
          )}
        </div>
        <div>
            <button onClick={handleSave} style={{ marginLeft: '8px', padding: '8px 16px', borderRadius: '4px', backgroundColor: '#007bff', color: '#ffffff', border: 'none', cursor: 'pointer' }}>
              Save
            </button>
        </div>
      </div>

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

          <button
            onClick={() => editor?.chain().focus().setImage({ src: 'https://i.natgeofe.com/n/548467d8-c5f1-4551-9f58-6817a8d2c45e/NationalGeographic_2572187_square.jpg' }).run()}
            className={editor?.isActive({ heading: { level: 1 } }) ? 'is-active' : ''}
          >
            Add Image
          </button>
        </div>
      </div>
      
      <div className="editor-container" style={{ marginBottom: '1rem' }}>
        <EditorContent editor={editor} className="editor-content" />
      </div>
      
    </div>
  )
}
