import React, { useState } from "react"
import {
  Modal, StyleSheet, Text, TextInput, TouchableOpacity, View
} from "react-native"
import { useStore } from "../lib/store"

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
      animationType="fade"
      transparent={true}
      visible={isVisible}
      onRequestClose={onClose}
    >
      <View style={styles.centeredView}>
        <View style={styles.modalView}>
          <Text style={styles.title}>Repository URL</Text>
          <TextInput
            style={styles.input}
            onChangeText={setNewRepoUrl}
            value={newRepoUrl}
            placeholder="Enter repository URL"
            placeholderTextColor="#888"
          />
          <View style={styles.buttonContainer}>
            <TouchableOpacity style={styles.button} onPress={handleSave}>
              <Text style={styles.buttonText}>Save</Text>
            </TouchableOpacity>
            <TouchableOpacity style={styles.button} onPress={onClose}>
              <Text style={styles.buttonText}>Cancel</Text>
            </TouchableOpacity>
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
    backgroundColor: '#000',
    borderRadius: 20,
    padding: 35,
    alignItems: 'center',
    shadowColor: '#fff',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.25,
    shadowRadius: 4,
    elevation: 5,
    borderColor: '#fff',
    borderWidth: 1,
  },
  title: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 15,
    color: '#fff',
    fontFamily: 'JetBrainsMono_400Regular',
  },
  input: {
    height: 40,
    width: 250,
    borderColor: '#fff',
    borderWidth: 1,
    marginBottom: 20,
    paddingHorizontal: 10,
    color: '#fff',
    fontFamily: 'JetBrainsMono_400Regular',
  },
  buttonContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    width: '100%',
  },
  button: {
    backgroundColor: '#333',
    padding: 10,
    borderRadius: 5,
    minWidth: 100,
    alignItems: 'center',
  },
  buttonText: {
    color: '#fff',
    fontFamily: 'JetBrainsMono_400Regular',
  },
});

export default RepoMenu;
