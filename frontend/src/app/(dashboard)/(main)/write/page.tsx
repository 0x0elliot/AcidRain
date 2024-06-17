"use client"

import { Button } from '@tremor/react';
import { siteConfig } from '@/app/siteConfig';
import { useEffect, useState } from 'react';
// get cookies from nookies
import cookies, { destroyCookie } from 'nookies';
import { Card } from '@tremor/react';

const PostCard = ({ post }) => {
  const characterCount = post.content.length;

  const formatDateTime = (dateTimeString) => {
    const dateTime = new Date(dateTimeString);
    return dateTime.toLocaleString(); // Adjust format as needed
  };

  return (
    <Card className="mx-auto max-w-xs mt-4 rounded-lg shadow-md">
      <div className="px-4 py-3">
        <p className="text-gray-800 text-lg font-semibold mb-2">{post.title}</p>
        <div className="flex justify-between mb-2">
          <p className="text-gray-600 text-sm">Last Updated: {formatDateTime(post.updated_at)}</p>
          <p className="text-gray-600 text-sm">Created At: {formatDateTime(post.created_at)}</p>
        </div>
        <p className="text-gray-700 text-sm mb-2">Character Count: {characterCount}</p>
        <a href={`/write/post?post_id=${post.id}`} className="text-blue-500 text-sm hover:underline">
          Read More
        </a>
      </div>
    </Card>
  );
};
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
        title: 'Your interesting title',
        content: '<p style="font-family: Trebuchet MS !important;">Start writing here...</p>',
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
            setPosts(data.posts);
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


      <div className="mt-6">
        <h2 className="text-lg font-semibold text-gray-900 sm:text-xl dark:text-gray-50">
          Your Posts
        </h2>
        <div className="mt-4 grid grid-cols-3 gap-4">
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </div>
      </div>

    </>
  )
}
