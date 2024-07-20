// this file will eventually belong in a CDN
document.addEventListener('DOMContentLoaded', function () {
  let baseUrl = window.location.origin;
  if (window.Shopify) {
    if (window.Shopify.shop) {
      baseUrl = "https://" + window.Shopify.shop;
    }
  }
  if ('serviceWorker' in navigator && 'PushManager' in window) {
    navigator.serviceWorker.register(`${baseUrl}/apps/acidrain/public/service-worker.js`)
      .then(function (registration) {
        console.log('Service Worker registered with scope:', registration.scope);

        // Check if the user is already subscribed
        registration.pushManager.getSubscription()
          .then(async function (subscription) {
            if (subscription) {
              console.log('Already subscribed:', subscription);

              // check local storage for subscription to avoid spam
              if (localStorage.getItem('acidRainWebPush') === '1') {
                console.log('Already subscribed:', subscription);
                return;
              }

              sendSubscriptionToServer(subscription);
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

  // tracking
  // const fpPromise = import('https://openfpcdn.io/fingerprintjs/v4')
  const fpPromise = import(`${baseUrl}/apps/acidrain/public/fingerprint.js`)
    .then(FingerprintJS => FingerprintJS.load())
  // Get the visitor identifier when you need it.
  fpPromise
    .then(fp => fp.get())
    .then(result => {
      // This is the visitor identifier:
      const visitorId = result.visitorId
      var shopifyUniqueId;
      localStorage.setItem('acidRainVisitorId', visitorId);

      try {
        shopifyUniqueId = window.ShopifyAnalytics.lib.user().traits().uniqToken;
      } catch (error) {
        console.error('Error getting shopify unique id:', error);
      }

      // sync visitor it all with the server
      fetch(`${baseUrl}/apps/acidrain/api/tracking/sync`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
          fingerprint: visitorId,
          store: baseUrl,
          shopify_unique_id: shopifyUniqueId,
          push_notification_subscription: localStorage.getItem('acidRainWebPush') === '1' ? localStorage.getItem('acidRainWebPushSubscription') : ""
        })
      })
        .then(function (response) {
          if (!response.ok) {
            throw new Error('Failed to sync visitor with server');
          }
          return response.json();
        })
        .then(function (data) {
          console.log('Visitor synced with server:', data);
        })
        .catch(function (error) {
          console.error('Error syncing visitor with server:', error);
        });

    })

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
  let storeUrl;
  storeUrl = window.location.origin;
  if (window.Shopify) {
    storeUrl = "https://" + window.Shopify.shop;
  }

  // get public key from the server
  fetch(`${storeUrl}/apps/acidrain/api/web-push-public-key`, {
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
    registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: applicationServerKey
    })
      .then(function (subscription) {
        console.log('User is subscribed:', subscription);
        sendSubscriptionToServer(subscription);
      })
      .catch(function (error) {
        console.error('Failed to subscribe the user:', error);
      });
  }).catch(function (error) {
    console.error('Error getting public key from server:', error);
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
  let storeUrl;
  requestObj = {
    subscription: subscription,
    storeUrl: storeUrl
  };

  let baseUrl = window.location.origin;
  if (window.Shopify) {
    baseUrl = "https://" + window.Shopify.shop;
  }

  fetch(`${baseUrl}/apps/acidrain/api/notification/subscribe`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(requestObj)
  })
    .then(function (response) {
      if (!response.ok) {
        throw new Error('Failed to send subscription to server');
      }
      return response.json();
    })
    .then(function (data) {
      localStorage.setItem('acidRainWebPush', '1');
      localStorage.setItem('acidRainWebPushSubscription', JSON.stringify(subscription));
      console.log('Subscription sent to server:', data);
    })
    .catch(function (error) {
      console.error('Error sending subscription to server:', error);
    });
}
