<?xml version="1.0" encoding="UTF-8" ?>
<lexer xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
     xsi:noNamespaceSchemaLocation="lexer.xsd" 
>
    <name>Date1</name>
    <author>Guillaume de Lestanville</author>
    <description>Analyse d'une date au format français</description>
    <dateCreation>2020-12-10</dateCreation>
    <initial>0</initial>
    <rules>
        <rule id="1" from="0" to="0" repeat="2">
            <test>
                <charset>DIGIT</charset>
            </test>
            <concat>true</concat>
        </rule>
        <rule id="2" from="0" to="1">
            <test>
                <in>-/.</in>
            </test>
            <reset>true</reset>
            <action>setJour</action>
        </rule>
        <rule id="3" from="0" to="3">
            <test><eos></eos></test>
            <action>setJour</action>
        </rule>
        <rule id="4" from="1" to="1" repeat="2">
            <test>
                <charset>DIGIT</charset>
            </test>
            <concat>true</concat>
        </rule>
        <rule id="5" from="1" to="2">
            <test>
                <in>-/.</in>
            </test>
            <reset>true</reset>
            <action>setMois</action>
        </rule>
        <rule id="6" from="1" to="3">
            <test><eos></eos></test>
            <action>setMois</action>
        </rule>
        <rule id="7" from="2" to="2" repeat="4">
            <test>
                <charset>DIGIT</charset>
            </test>
            <concat>true</concat>
        </rule>
        <rule id="8" from="2" to="3">
            <test><eos></eos></test>
            <action>setAnnee</action>
        </rule>
        <rule id="9" from="3" to="3">
            <test><any/></test>
            <final>true</final>
        </rule>
    </rules>
</lexer>
