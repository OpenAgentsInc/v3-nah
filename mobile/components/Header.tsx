import React from 'react';
import { View, Image, StyleSheet } from 'react-native';
import RelayStatusIcon from './RelayStatusIcon';

interface HeaderProps {
  isConnected: boolean;
}

const Header: React.FC<HeaderProps> = ({ isConnected }) => {
  return (
    <View style={styles.header}>
      <Image source={require('../assets/sqlogo-t.png')} style={styles.logo} resizeMode="contain" />
      <RelayStatusIcon isConnected={isConnected} />
    </View>
  );
};

const styles = StyleSheet.create({
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingTop: 0,
    paddingHorizontal: 20,
    height: 50,
  },
  logo: {
    width: 40,
    height: 40,
  },
});

export default Header;