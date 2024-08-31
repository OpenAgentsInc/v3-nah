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
    width: 10,
    height: 10,
    borderRadius: 5,
  },
  connected: {
    backgroundColor: '#4CAF50',
  },
  disconnected: {
    backgroundColor: '#F44336',
  },
});

export default RelayStatusIcon;