import React from 'react';
import { View, StyleSheet } from 'react-native';

interface RelayStatusIconProps {
  isConnected: boolean;
}

const RelayStatusIcon: React.FC<RelayStatusIconProps> = ({ isConnected }) => {
  return (
    <View style={[styles.icon, isConnected ? styles.connected : styles.disconnected]} />
  );
};

const styles = StyleSheet.create({
  icon: {
    width: 12,
    height: 12,
    borderRadius: 6,
    position: 'absolute',
    top: 40,
    right: 20,
  },
  connected: {
    backgroundColor: '#4CAF50',
  },
  disconnected: {
    backgroundColor: '#F44336',
  },
});

export default RelayStatusIcon;