"use client"

import { Button } from '@tremor/react';
import { siteConfig } from '@/app/siteConfig';
import { useEffect, useState } from 'react';
// get cookies from nookies
import cookies, { destroyCookie } from 'nookies';

export default function Write() {
  const [posts, setPosts] = useState([]);

  async function handleStartWriting() {
    let accessToken = cookies.get(null).access_token;

    let r = fetch(`${siteConfig.baseApiUrl}/api/post/private/set`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': `Bearer ${accessToken}`,
      },
      mode: 'cors',
      body: JSON.stringify({
        title: 'New Post',
        content: 'Start writing your post here.',
      }),
    })

    r.then((response) => {
      if (response.ok) {
        // redirect to the new post
        response.json().then((data) => {
          let post = data.post;

          window.location.href = `/write/post?post_id=${post.id}`;
        });
      } else if (response.status === 401) {

        destroyCookie(null, 'access_token');
        destroyCookie(null, 'refresh_token');

        window.location.href = '/login';
      }
    });

  }


  useEffect(() => {
    // check accessToken from cookies
    const accessToken_ = cookies.get(null).access_token;

    // Fetch all posts
    fetch(`${siteConfig.baseApiUrl}/api/post/private/all`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        "Authorization": `Bearer ${accessToken_}`,
      },
    })
      .then((response) => {
        if (response.ok) {
          response.json().then((data) => {
            setPosts(data);
          });
        } else if (response.status === 401) {
          window.location.href = '/login';
        }
      });
    
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
