import { useState, useEffect } from 'react';
import { Audio } from 'expo-av';

export const useAudioPermission = () => {
  const [permissionStatus, setPermissionStatus] = useState<'granted' | 'denied' | 'pending'>('pending');

  useEffect(() => {
    (async () => {
      const { status } = await Audio.requestPermissionsAsync();
      setPermissionStatus(status === 'granted' ? 'granted' : 'denied');
    })();
  }, []);

  return permissionStatus;
};