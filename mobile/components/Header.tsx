import React, { useState, useMemo } from 'react';
import { View, Image, StyleSheet, TouchableOpacity, Text } from 'react-native';
import RelayStatusIcon from './RelayStatusIcon';
import RepoMenu from './RepoMenu';
import { useStore } from '../lib/store';

interface HeaderProps {
  isConnected: boolean;
}

const Header: React.FC<HeaderProps> = ({ isConnected }) => {
  const [isMenuVisible, setIsMenuVisible] = useState(false);
  const { activeRepoUrl } = useStore();

  const handleLogoPress = () => {
    setIsMenuVisible(true);
  };

  const handleCloseMenu = () => {
    setIsMenuVisible(false);
  };

  const displayRepoUrl = useMemo(() => {
    const url = new URL(activeRepoUrl);
    const pathParts = url.pathname.split('/').filter(Boolean);
    if (pathParts.length >= 2) {
      return `${pathParts[0]}/${pathParts[1]}`;
    }
    return activeRepoUrl;
  }, [activeRepoUrl]);

  return (
    <View style={styles.header}>
      <TouchableOpacity onPress={handleLogoPress}>
        <Image source={require('../assets/sqlogo-t.png')} style={styles.logo} resizeMode="contain" />
      </TouchableOpacity>
      <View style={styles.repoUrlContainer}>
        <Text style={styles.repoUrlText} numberOfLines={1} ellipsizeMode="middle">
          {displayRepoUrl}
        </Text>
      </View>
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
    backgroundColor: '#000',
    borderBottomColor: 'rgba(255, 255, 255, 0.5)',
    borderBottomWidth: 1,
  },
  logo: {
    width: 40,
    height: 40,
  },
  repoUrlContainer: {
    flex: 1,
    marginHorizontal: 10,
  },
  repoUrlText: {
    color: 'rgba(255, 255, 255, 0.5)',
    fontFamily: 'JetBrainsMono',
    fontSize: 12,
  },
});

export default Header;