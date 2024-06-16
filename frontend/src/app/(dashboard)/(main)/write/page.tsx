"use client"

import { Button } from '@tremor/react';
import { siteConfig } from '@/app/siteConfig';
import { useEffect } from 'react';

export default function Write() {
  async function handleStartWriting() {
    // create a new post
    fetch(`${siteConfig.baseApiUrl}/api/post/private/set`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({
        title: 'New Post',
        content: 'Start writing your post here.',
      }),
    }).then((response) => {
      if (response.ok) {
        // redirect to the new post
        response.json().then((data) => {
          window.location.href = `/write/post/${data.id}`;
        });
      } else {
        console.error('Failed to create a new post.');
      }
    });
  }

  useEffect(() => {
    const startWritingButton = document.getElementById('startWritingButton');
    if (startWritingButton) {
      startWritingButton.addEventListener('click', handleStartWriting);
    }

    return () => {
      if (startWritingButton) {
        startWritingButton.removeEventListener('click', handleStartWriting);
      }
    };
  }, []);

  return (
    <>
      <h1 className="text-lg font-semibold text-gray-900 sm:text-xl dark:text-gray-50">
        Manage Posts
      </h1>

      <div className="mt-6">
        <Button id="startWritingButton">
          New Post
        </Button>
      </div>

    </>
  )
}
