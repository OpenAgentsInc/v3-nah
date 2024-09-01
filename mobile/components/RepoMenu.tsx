import React, { useState } from 'react';
import { View, TextInput, Button, StyleSheet, Modal } from 'react-native';
import { useStore } from '../lib/store';

interface RepoMenuProps {
  isVisible: boolean;
  onClose: () => void;
}

const RepoMenu: React.FC<RepoMenuProps> = ({ isVisible, onClose }) => {
  const { activeRepoUrl, setActiveRepoUrl } = useStore();
  const [newRepoUrl, setNewRepoUrl] = useState(activeRepoUrl);

  const handleSave = () => {
    setActiveRepoUrl(newRepoUrl);
    onClose();
  };

  return (
    <Modal
      animationType="slide"
      transparent={true}
      visible={isVisible}
      onRequestClose={onClose}
    >
      <View style={styles.centeredView}>
        <View style={styles.modalView}>
          <TextInput
            style={styles.input}
            onChangeText={setNewRepoUrl}
            value={newRepoUrl}
            placeholder="Enter repository URL"
          />
          <View style={styles.buttonContainer}>
            <Button title="Save" onPress={handleSave} />
            <Button title="Cancel" onPress={onClose} />
          </View>
        </View>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  centeredView: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
  },
  modalView: {
    backgroundColor: 'white',
    borderRadius: 20,
    padding: 35,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.25,
    shadowRadius: 4,
    elevation: 5,
  },
  input: {
    height: 40,
    width: 250,
    borderColor: 'gray',
    borderWidth: 1,
    marginBottom: 20,
    paddingHorizontal: 10,
  },
  buttonContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    width: '100%',
  },
});

export default RepoMenu;