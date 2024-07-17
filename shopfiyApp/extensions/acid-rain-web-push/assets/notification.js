// this file will eventually belong in a CDN
document.addEventListener('DOMContentLoaded', function () {
    if ('serviceWorker' in navigator && 'PushManager' in window) {
      navigator.serviceWorker.register('/apps/acidrain/public/service-worker.js')
        .then(function (registration) {
          console.log('Service Worker registered with scope:', registration.scope);
  
          // Check if the user is already subscribed
          registration.pushManager.getSubscription()
            .then(function (subscription) {
              if (subscription) {
                console.log('Already subscribed:', subscription);
                // You can skip asking for permission if already subscribed
                sendSubscriptionToServer(subscription);
                requestTestNotification(subscription);
              } else {
                askForNotificationPermission(registration);
              }
            });
        })
        .catch(function (error) {
          console.error('Service Worker registration failed:', error);
        });
    } else {
      console.warn('Push messaging is not supported');
    }
  });
  
  function askForNotificationPermission(registration) {
    Notification.requestPermission().then(function (permission) {
      if (permission === 'granted') {
        console.log('Notification permission granted.');
        subscribeUserToPush(registration);
      } else {
        console.warn('Notification permission denied.');
      }
    });
  }
  
  function subscribeUserToPush(registration) {
    let applicationServerKey;

    // get public key from the server
    fetch('/apps/acidrain/api/web-push-public-key', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    }).then(function (response) {
      if (!response.ok) {
        throw new Error('Failed to get public key from server');
      }
      return response.json();
    }).then(function (data) {
      applicationServerKey = urlB64ToUint8Array(data.publicKey);
    }).catch(function (error) {
      console.error('Error getting public key from server:', error);
    });

    registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: applicationServerKey
    })
    .then(function (subscription) {
      console.log('User is subscribed:', subscription);
      sendSubscriptionToServer(subscription);
      requestTestNotification(subscription);
    })
    .catch(function (error) {
      console.error('Failed to subscribe the user:', error);
    });
  }
  
  function urlB64ToUint8Array(base64String) {
    const padding = '='.repeat((4 - base64String.length % 4) % 4);
    const base64 = (base64String + padding)
      .replace(/-/g, '+')
      .replace(/_/g, '/');
  
    const rawData = window.atob(base64);
    const outputArray = new Uint8Array(rawData.length);
  
    for (let i = 0; i < rawData.length; ++i) {
      outputArray[i] = rawData.charCodeAt(i);
    }
    return outputArray;
  }
  
  function sendSubscriptionToServer(subscription) {
    fetch('apps/acidrain/api/notification/subscribe', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(subscription)
    })
    .then(function (response) {
      if (!response.ok) {
        throw new Error('Failed to send subscription to server');
      }
      return response.json();
    })
    .then(function (data) {
      console.log('Subscription sent to server:', data);
    })
    .catch(function (error) {
      console.error('Error sending subscription to server:', error);
    });
  }
  
  function requestTestNotification(subscription) {
    fetch('http://localhost:3000/test-notification', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(subscription)
    })
    .then(function (response) {
      if (!response.ok) {
        throw new Error('Failed to request test notification');
      }
      return response.json();
    })
    .then(function (data) {
      console.log('Test notification requested:', data);
    })
    .catch(function (error) {
      console.error('Error requesting test notification:', error);
    });
  }
  