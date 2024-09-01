import React, { useState } from 'react';
import { View, Image, StyleSheet, TouchableOpacity } from 'react-native';
import RelayStatusIcon from './RelayStatusIcon';
import RepoMenu from './RepoMenu';

interface HeaderProps {
  isConnected: boolean;
}

const Header: React.FC<HeaderProps> = ({ isConnected }) => {
  const [isMenuVisible, setIsMenuVisible] = useState(false);

  const handleLogoPress = () => {
    setIsMenuVisible(true);
  };

  const handleCloseMenu = () => {
    setIsMenuVisible(false);
  };

  return (
    <View style={styles.header}>
      <TouchableOpacity onPress={handleLogoPress}>
        <Image source={require('../assets/sqlogo-t.png')} style={styles.logo} resizeMode="contain" />
      </TouchableOpacity>
      <RelayStatusIcon isConnected={isConnected} />
      <RepoMenu isVisible={isMenuVisible} onClose={handleCloseMenu} />
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