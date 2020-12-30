import React, { useState, useEffect } from 'react';
import { View,
    AsyncStorage,
    KeyboardAvoidingView,
    BackHandler,
    Alert,
    Image,
    Text,
    TextInput,
    TouchableOpacity,
    StyleSheet} from 'react-native';
import logo from '../assets/logo.png';
import { Platform } from '@unimodules/core';

import api from '../services/api'
export default function User({navigation}) {
    const [name, setName] = useState('');
    const [doc,   setDoc] = useState('');
    const [pass, setPass] = useState('');

    const [phone, setPhon] = useState('');
    const [token, setToke] = useState('');

    const [user, setUser] = useState(null);

    function handleBackButtonClick() {
        console.log("backing...");
        navigation.navigate(navigation.getParam('back'));
        return true;
    }

    useEffect(() => {
        BackHandler.addEventListener('hardwareBackPress', handleBackButtonClick);
        return () => {
            BackHandler.removeEventListener('hardwareBackPress', handleBackButtonClick);
        };
    }, []);


    useEffect(()=>{
        const luser = navigation.getParam('user');
        AsyncStorage.getItem('curr_user').then(curr=>{
            if (curr) {
                curr = JSON.parse(curr);
                setToke(curr.token);
            }
        });

        if(luser) {
            setName(luser.name);
            setPhon(luser.phone);
            setDoc(luser.document);

            setUser(luser);
        }
    },[]);

    async function handleSubmit() {
        var url = '/user';
        if (user) {
            url += `/${user.id}`;
            setPass(null);
        }

        console.log(url,' ',token);
        try {
            const response = await api.post(url, {token,data: {
                password: pass,
                document: doc,
                phone,
                name,
            }})

            navigation.navigate('List');
        } catch {
            Alert.alert(`Não foi possível ${user?'editar':'criar'} ${name}`);
        }

    }

    return (
        <KeyboardAvoidingView
         enabled={Platform.OS== 'ios'} 
         behavior="padding" style={styles.container}>
            <View style={styles.box}>
                <Image source={logo}/>

                <View style={styles.form}>

                    <Text style={styles.label}>Nome Completo</Text>
                    <TextInput
                        style={styles.input}
                        placeholder="Digite seu nome completo"
                        placeholderTextColor="#999"
                        keyboardType="default"
                        autoCapitalize="words"
                        autoCorrect={false}
                        onChangeText={setName}
                    >{name}</TextInput>

                    <Text style={styles.label}>Número de Telefone</Text>
                    <TextInput
                        style={styles.input}
                        placeholder="Digite seu telefone"
                        placeholderTextColor="#999"
                        keyboardType="default"
                        autoCapitalize="words"
                        autoCorrect={false}
                        onChangeText={setPhon}
                    >{phone}</TextInput>

                    <Text style={styles.label}>CPF</Text>
                    <TextInput
                        style={styles.input}
                        placeholder="Digite seu CPF"
                        placeholderTextColor="#999"
                        keyboardType="email-address"
                        autoCapitalize="none"
                        autoCorrect={false}
                        value={doc}
                        onChangeText={setDoc}
                    >
                    </TextInput>

                    {!user?(<><Text style={styles.label}>SENHA</Text>
                    <TextInput
                        style={styles.input}
                        placeholder="Digite sua senha"
                        placeholderTextColor="#999"
                        keyboardType={"default"}
                        secureTextEntry={true}
                        onChangeText={setPass}
                    >{pass}</TextInput></>):null}

                    <TouchableOpacity onPress={handleSubmit} style={styles.button}>
                        <Text style={styles.buttonText}>{user?'Salvar':'Criar usuário'}</Text>
                    </TouchableOpacity>
                </View>

            </View>
        </KeyboardAvoidingView>
    )
}
const styles = StyleSheet.create({
   container: {
       flex:1,
       justifyContent:'center',
       alignItems:'center'
   }, 
   box:{
        alignSelf: 'stretch',
        paddingHorizontal: 30,
        alignItems: 'center',
        justifyContent: 'space-between',
        marginTop: 30,
   },
   form:{
        alignSelf: 'stretch',
        paddingHorizontal: 30,
        marginTop: 30,
        marginBottom: 30,
   },
   label: {
       fontWeight: 'bold',
       color:'#444',
       marginBottom:8
   },
   input: {
       borderWidth:1,
       borderColor: '#ddd',
       paddingHorizontal:20,
       fontSize: 16,
       color:'#444',
       height: 44,
       marginBottom: 20,
       borderRadius: 2
   },
   button: {
       height: 42,
       backgroundColor: '#23B185',
       justifyContent: 'center',
       alignItems:'center',
       borderRadius:2,
   },
   buttonText:{
       color: '#FFF',
       fontWeight:'bold',
       fontSize:16,
   },
   linkText: {
       color: '#23B185',
       fontWeight:'bold',
       fontSize:16,
   },
});